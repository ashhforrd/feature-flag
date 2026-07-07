import test from "node:test"
import assert from "node:assert/strict"
import { FeatureFlagClient } from "../src/index.js"

test("requires baseUrl", () => {
  assert.throws(() => {
    new FeatureFlagClient()
  }, /baseUrl is required/)
})

test("removes trailing slash from baseUrl", () => {
  const client = new FeatureFlagClient({
    baseUrl: "http://localhost:8080/"
  })

  assert.equal(client.baseUrl, "http://localhost:8080")
})

test("evaluate returns backend response", async () => {
  const originalFetch = globalThis.fetch

  globalThis.fetch = async (url, options) => {
    assert.equal(url, "http://localhost:8080/flags/new-checkout/evaluate")
    assert.equal(options.method, "POST")
    assert.equal(options.headers["Content-Type"], "application/json")

    assert.deepEqual(JSON.parse(options.body), {
      user: {
        id: "user-123",
        attributes: {
          country: "ID"
        }
      },
      defaultValue: false
    })

    return {
      ok: true,
      async json() {
        return {
          flagKey: "new-checkout",
          enabled: true,
          reason: "DEFAULT_RULE"
        }
      }
    }
  }

  try {
    const client = new FeatureFlagClient({
      baseUrl: "http://localhost:8080"
    })

    const result = await client.evaluate(
      "new-checkout",
      {
        id: "user-123",
        attributes: {
          country: "ID"
        }
      },
      false
    )

    assert.deepEqual(result, {
      flagKey: "new-checkout",
      enabled: true,
      reason: "DEFAULT_RULE"
    })
  } finally {
    globalThis.fetch = originalFetch
  }
})

test("isEnabled returns only enabled boolean", async () => {
  const originalFetch = globalThis.fetch

  globalThis.fetch = async () => {
    return {
      ok: true,
      async json() {
        return {
          flagKey: "new-checkout",
          enabled: true,
          reason: "DEFAULT_RULE"
        }
      }
    }
  }

  try {
    const client = new FeatureFlagClient({
      baseUrl: "http://localhost:8080"
    })

    const enabled = await client.isEnabled("new-checkout", { id: "user-123" }, false)

    assert.equal(enabled, true)
  } finally {
    globalThis.fetch = originalFetch
  }
})

test("evaluate returns default value when request fails", async () => {
  const originalFetch = globalThis.fetch

  globalThis.fetch = async () => {
    throw new Error("network error")
  }

  try {
    const client = new FeatureFlagClient({
      baseUrl: "http://localhost:8080"
    })

    const result = await client.evaluate("new-checkout", { id: "user-123" }, true)

    assert.deepEqual(result, {
      flagKey: "new-checkout",
      enabled: true,
      reason: "SDK_REQUEST_FAILED"
    })
  } finally {
    globalThis.fetch = originalFetch
  }
})