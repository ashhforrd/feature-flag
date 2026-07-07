export class FeatureFlagClient {
    constructor(options) {
        if (!options || !options.baseUrl) {
            throw new Error("baseUrl is required")
        }

        this.baseUrl = options.baseUrl.replace(/\/$/, "")
    }

    async evaluate(flagKey, user, defaultValue = false) {
        try {
            const response = await fetch(`${this.baseUrl}/flags/${flagKey}/evaluate`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    user,
                    defaultValue
                })
            })

            if (!response.ok) {
                return {
                    flagKey,
                    enabled: defaultValue,
                    reason: "SDK_REQUEST_FAILED"
                }
            }

            return await response.json()
        } catch (err) {
            return {
                flagKey,
                enabled: defaultValue,
                reason: "SDK_REQUEST_FAILED"
            }
        }
    }

    async isEnabled(flagKey, user, defaultValue = false) {
        const result = await this.evaluate(flagKey, user, defaultValue)
        return result.enabled
    }
}