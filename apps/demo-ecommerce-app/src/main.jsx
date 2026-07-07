import React, { useEffect, useMemo, useState } from "react"
import { createRoot } from "react-dom/client"
import { FeatureFlagClient } from "../../../packages/js-sdk/src/index.js"
import "./styles.css"

const client = new FeatureFlagClient({
  baseUrl: import.meta.env.VITE_CONFIG_SERVICE_URL || "/api"
})

const users = [
  {
    id: "user-123",
    name: "Alice Tan",
    segment: "Regular customer",
    attributes: {
      country: "ID",
      email: "alice@example.com"
    }
  },
  {
    id: "internal-001",
    name: "Raka Internal",
    segment: "Internal tester",
    attributes: {
      country: "ID",
      email: "raka@company.test"
    }
  },
  {
    id: "user-456",
    name: "Maya Lee",
    segment: "International customer",
    attributes: {
      country: "SG",
      email: "maya@example.com"
    }
  }
]

const cartItems = [
  { name: "Everyday Linen Shirt", color: "Olive", price: 68, swatch: "from-stone-300 to-olive-700" },
  { name: "Canvas Market Tote", color: "Natural", price: 42, swatch: "from-amber-100 to-stone-300" },
  { name: "Matte Steel Bottle", color: "Graphite", price: 35, swatch: "from-zinc-400 to-zinc-800" }
]

function formatCurrency(value) {
  return new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD"
  }).format(value)
}

function App() {
  const [selectedUserId, setSelectedUserId] = useState(users[0].id)
  const [evaluation, setEvaluation] = useState(null)
  const [status, setStatus] = useState("loading")
  const [conversionStatus, setConversionStatus] = useState("idle")

  const user = useMemo(
    () => users.find((item) => item.id === selectedUserId) || users[0],
    [selectedUserId]
  )

  const subtotal = cartItems.reduce((sum, item) => sum + item.price, 0)
  const shipping = evaluation?.enabled ? 0 : 8
  const total = subtotal + shipping

  useEffect(() => {
    let cancelled = false

    async function loadFlag() {
      setStatus("loading")
      setConversionStatus("idle")
      const result = await client.evaluate("new-checkout", user, false)

      if (!cancelled) {
        setEvaluation(result)
        setStatus("ready")
      }
    }

    loadFlag()

    return () => {
      cancelled = true
    }
  }, [user])

  async function handleCheckoutComplete() {
    setConversionStatus("recording")

    const recorded = await client.recordConversion(
      "new-checkout",
      user.id,
      "checkout_completed"
    )

    setConversionStatus(recorded ? "recorded" : "failed")
  }

  const newCheckoutEnabled = evaluation?.enabled === true

  return (
    <main className="min-h-screen bg-[#f5f4ef] bg-[linear-gradient(90deg,rgba(24,27,31,0.04)_1px,transparent_1px),linear-gradient(180deg,rgba(24,27,31,0.04)_1px,transparent_1px)] bg-[size:48px_48px] px-4 py-4 text-[#181b1f] sm:px-8 sm:py-8">
      <section className="mx-auto w-full max-w-6xl">
        <header className="flex flex-col gap-6 rounded-t-lg border border-[#dedbd2] bg-[#fffefc]/90 p-7 shadow-[0_18px_50px_rgba(31,34,38,0.08)] sm:flex-row sm:items-center sm:justify-between">
          <div>
            <p className="mb-2 text-xs font-bold uppercase text-[#6f5f4d]">Northstar Supply</p>
            <h1 className="text-5xl font-bold leading-none tracking-normal sm:text-7xl">Checkout</h1>
          </div>
          <div
            className={`w-fit min-w-40 rounded-full px-4 py-2.5 text-center text-sm font-bold ${
              newCheckoutEnabled ? "bg-emerald-100 text-emerald-950" : "bg-orange-100 text-orange-950"
            }`}
          >
            {newCheckoutEnabled ? "New checkout" : "Classic checkout"}
          </div>
        </header>

        <section className="grid gap-4 border-x border-b border-[#dedbd2] bg-[#fffefc]/90 p-5 shadow-[0_18px_50px_rgba(31,34,38,0.08)] md:grid-cols-[auto_minmax(220px,1fr)_auto] md:items-center">
          <label className="text-sm font-bold text-zinc-600" htmlFor="user-select">
            Customer
          </label>
          <select
            className="min-h-11 w-full rounded-md border border-stone-300 bg-[#fffefa] px-3 text-sm text-zinc-950 outline-none focus:border-zinc-800"
            id="user-select"
            value={selectedUserId}
            onChange={(event) => setSelectedUserId(event.target.value)}
          >
            {users.map((item) => (
              <option key={item.id} value={item.id}>
                {item.name} - {item.segment}
              </option>
            ))}
          </select>
          <div className="flex flex-wrap gap-2 md:justify-end">
            <span className="rounded-full border border-stone-300 bg-stone-50 px-3 py-1.5 text-xs font-medium text-zinc-600">
              {status === "loading" ? "Evaluating" : evaluation?.reason}
            </span>
            {evaluation?.bucket !== undefined && (
              <span className="rounded-full border border-stone-300 bg-stone-50 px-3 py-1.5 text-xs font-medium text-zinc-600">
                Bucket {evaluation.bucket}
              </span>
            )}
          </div>
        </section>

        <div className="mt-5 grid gap-5 lg:grid-cols-[minmax(0,1.35fr)_minmax(320px,0.65fr)]">
          {newCheckoutEnabled ? (
            <NewCheckout
              user={user}
              shipping={shipping}
              total={total}
              conversionStatus={conversionStatus}
              onCheckoutComplete={handleCheckoutComplete}
            />
          ) : (
            <ClassicCheckout
              user={user}
              total={total}
              conversionStatus={conversionStatus}
              onCheckoutComplete={handleCheckoutComplete}
            />
          )}

          <OrderSummary subtotal={subtotal} shipping={shipping} total={total} />
        </div>
      </section>
    </main>
  )
}

function ClassicCheckout({ user, total, conversionStatus, onCheckoutComplete }) {
  return (
    <section className="rounded-lg border border-[#dedbd2] bg-[#fffefc]/90 p-7 shadow-[0_18px_50px_rgba(31,34,38,0.08)]">
      <div className="mb-7">
        <p className="mb-2 text-xs font-bold uppercase text-[#6f5f4d]">Stable flow</p>
        <h2 className="text-3xl font-bold tracking-normal">Classic checkout</h2>
      </div>

      <div className="grid gap-4">
        <ReadOnlyField label="Email" value={user.attributes.email} />
        <ReadOnlyField label="Shipping country" value={user.attributes.country} />
        <ReadOnlyField label="Payment method" value="Visa ending in 4242" />
      </div>

      <button
        className="mt-6 min-h-12 w-full rounded-md bg-zinc-950 px-5 font-bold text-white disabled:cursor-not-allowed disabled:bg-zinc-400"
        disabled={conversionStatus === "recording"}
        onClick={onCheckoutComplete}
      >
        {conversionStatus === "recording" ? "Recording order" : `Place order ${formatCurrency(total)}`}
      </button>
      <ConversionStatus status={conversionStatus} />
    </section>
  )
}

function NewCheckout({ user, shipping, total, conversionStatus, onCheckoutComplete }) {
  return (
    <section className="rounded-lg border border-emerald-200 bg-[#fffefc]/90 p-7 shadow-[0_18px_50px_rgba(31,34,38,0.08)]">
      <div className="mb-7">
        <p className="mb-2 text-xs font-bold uppercase text-[#6f5f4d]">Flagged experience</p>
        <h2 className="text-3xl font-bold tracking-normal">One-page checkout</h2>
      </div>

      <div className="mb-5 grid gap-2 sm:grid-cols-3">
        {["Apple Pay", "Shop Pay", "Card"].map((method) => (
          <button key={method} className="min-h-12 rounded-md bg-emerald-100 px-4 font-bold text-emerald-950">
            {method}
          </button>
        ))}
      </div>

      <div className="grid gap-3">
        <ReviewRow label="Customer" value={user.name} />
        <ReviewRow label="Delivery" value={shipping === 0 ? "Free priority" : "Standard"} />
        <ReviewRow label="Due today" value={formatCurrency(total)} />
      </div>

      <button
        className="mt-6 min-h-12 w-full rounded-md bg-emerald-800 px-5 font-bold text-white disabled:cursor-not-allowed disabled:bg-emerald-300"
        disabled={conversionStatus === "recording"}
        onClick={onCheckoutComplete}
      >
        {conversionStatus === "recording" ? "Recording checkout" : "Complete secure checkout"}
      </button>
      <ConversionStatus status={conversionStatus} />
    </section>
  )
}

function ConversionStatus({ status }) {
  if (status === "idle") {
    return null
  }

  const styles = {
    recording: "border-stone-300 bg-stone-50 text-zinc-600",
    recorded: "border-emerald-200 bg-emerald-50 text-emerald-900",
    failed: "border-red-200 bg-red-50 text-red-900"
  }

  const messages = {
    recording: "Recording conversion event",
    recorded: "Conversion event recorded",
    failed: "Conversion event failed"
  }

  return (
    <p className={`mt-4 rounded-md border px-3 py-2 text-sm font-bold ${styles[status]}`}>
      {messages[status]}
    </p>
  )
}

function ReadOnlyField({ label, value }) {
  return (
    <label className="grid gap-2 text-sm font-bold text-zinc-600">
      {label}
      <input
        className="min-h-11 w-full rounded-md border border-stone-300 bg-[#fffefa] px-3 text-zinc-950"
        value={value}
        readOnly
      />
    </label>
  )
}

function ReviewRow({ label, value }) {
  return (
    <div className="flex min-h-14 items-center justify-between gap-4 rounded-md border border-emerald-100 bg-emerald-50/40 px-4 py-3">
      <span className="text-sm text-zinc-600">{label}</span>
      <strong className="text-right text-zinc-950">{value}</strong>
    </div>
  )
}

function OrderSummary({ subtotal, shipping, total }) {
  return (
    <aside className="rounded-lg border border-[#dedbd2] bg-[#fffefc]/90 p-7 shadow-[0_18px_50px_rgba(31,34,38,0.08)]">
      <h2 className="text-2xl font-bold tracking-normal">Order summary</h2>
      <div className="my-6 grid gap-4">
        {cartItems.map((item) => (
          <div className="grid grid-cols-[54px_1fr_auto] items-center gap-3" key={item.name}>
            <div className={`aspect-square w-[54px] rounded-lg bg-gradient-to-br ${item.swatch}`} aria-hidden="true" />
            <div>
              <strong className="block text-sm text-zinc-950">{item.name}</strong>
              <span className="text-sm text-zinc-500">{item.color}</span>
            </div>
            <p className="m-0 font-bold">{formatCurrency(item.price)}</p>
          </div>
        ))}
      </div>

      <div className="grid gap-3 border-t border-stone-200 pt-5">
        <TotalRow label="Subtotal" value={formatCurrency(subtotal)} />
        <TotalRow label="Shipping" value={shipping === 0 ? "Free" : formatCurrency(shipping)} />
        <TotalRow label="Total" value={formatCurrency(total)} large />
      </div>
    </aside>
  )
}

function TotalRow({ label, value, large = false }) {
  return (
    <div className={`flex items-center justify-between gap-4 ${large ? "text-xl" : "text-base"}`}>
      <span className="text-zinc-600">{label}</span>
      <strong className="text-zinc-950">{value}</strong>
    </div>
  )
}

createRoot(document.getElementById("root")).render(<App />)
