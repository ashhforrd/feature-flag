import React, { useEffect, useMemo, useState } from "react"
import { createRoot } from "react-dom/client"
import {
  CheckCircle2,
  ChevronDown,
  CreditCard,
  Package,
  RefreshCw,
  ShieldCheck,
  ShoppingBag,
  Truck,
  UserRound
} from "lucide-react"
import { FeatureFlagClient } from "@ashhforrd/feature-flags-js"
import "./styles.css"

const client = new FeatureFlagClient({
  baseUrl: import.meta.env.VITE_CONFIG_SERVICE_URL || "/api"
})

const users = [
  {
    id: "user-123",
    name: "Alice Tan",
    segment: "Regular customer",
    attributes: { country: "ID", email: "alice@example.com" }
  },
  {
    id: "internal-001",
    name: "Raka Internal",
    segment: "Internal tester",
    attributes: { country: "ID", email: "raka@company.test" }
  },
  {
    id: "user-456",
    name: "Maya Lee",
    segment: "International customer",
    attributes: { country: "SG", email: "maya@example.com" }
  }
]

const cartItems = [
  { name: "Everyday Linen Shirt", color: "Oxford blue", price: 68, swatch: "from-slate-200 to-blue-700" },
  { name: "Canvas Market Tote", color: "Mist", price: 42, swatch: "from-slate-100 to-slate-300" },
  { name: "Matte Steel Bottle", color: "Navy", price: 35, swatch: "from-blue-300 to-slate-900" }
]

function formatCurrency(value) {
  return new Intl.NumberFormat("en-US", { style: "currency", currency: "USD" }).format(value)
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
  const newCheckoutEnabled = evaluation?.enabled === true

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
    const recorded = await client.recordConversion("new-checkout", user.id, "checkout_completed")
    setConversionStatus(recorded ? "recorded" : "failed")
  }

  return (
    <main className="min-h-screen bg-slate-50 px-4 py-6 text-slate-950 sm:px-6 lg:px-8">
      <section className="mx-auto w-full max-w-6xl">
        <header className="rounded-lg border border-slate-200 bg-white px-6 py-5">
          <div className="flex flex-col gap-5 md:flex-row md:items-center md:justify-between">
            <div>
              <p className="flex items-center gap-2 text-xs font-semibold uppercase text-slate-500">
                <ShoppingBag className="h-3.5 w-3.5" />
                Northstar Supply
              </p>
              <h1 className="mt-1 text-3xl font-semibold tracking-normal text-slate-950 md:text-4xl">Checkout</h1>
            </div>
            <div
              className={`inline-flex w-fit items-center gap-2 rounded-full border px-3 py-1.5 text-sm font-semibold ${
                newCheckoutEnabled
                  ? "border-blue-200 bg-blue-50 text-blue-700"
                  : "border-slate-200 bg-slate-100 text-slate-700"
              }`}
            >
              <ShieldCheck className="h-4 w-4" />
              {newCheckoutEnabled ? "New checkout" : "Classic checkout"}
            </div>
          </div>

          <div className="mt-6 grid gap-3 border-t border-slate-100 pt-5 md:grid-cols-[120px_minmax(220px,1fr)_auto] md:items-center">
            <label className="text-sm font-medium text-slate-600" htmlFor="user-select">
              Customer
            </label>
            <div className="relative">
              <select
                className="min-h-10 w-full appearance-none rounded-md border border-slate-200 bg-white px-3 pr-9 text-sm text-slate-950 outline-none transition focus:border-blue-500 focus:ring-2 focus:ring-blue-100"
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
              <ChevronDown className="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400" />
            </div>
            <div className="flex flex-wrap gap-2 md:justify-end">
              <StatusPill icon={RefreshCw}>{status === "loading" ? "Evaluating" : evaluation?.reason}</StatusPill>
              {evaluation?.bucket !== undefined && <StatusPill icon={ShieldCheck}>Bucket {evaluation.bucket}</StatusPill>}
            </div>
          </div>
        </header>

        <div className="mt-5 grid gap-5 lg:grid-cols-[minmax(0,1fr)_360px]">
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

function StatusPill({ children, icon: Icon }) {
  return (
    <span className="inline-flex items-center gap-1.5 rounded-full border border-slate-200 bg-slate-50 px-3 py-1 text-xs font-medium text-slate-600">
      {Icon && <Icon className="h-3.5 w-3.5" />}
      {children}
    </span>
  )
}

function ClassicCheckout({ user, total, conversionStatus, onCheckoutComplete }) {
  return (
    <section className="rounded-lg border border-slate-200 bg-white p-6">
      <SectionHeading eyebrow="Stable flow" title="Classic checkout" icon={CreditCard} />
      <div className="grid gap-4">
        <ReadOnlyField label="Email" value={user.attributes.email} />
        <ReadOnlyField label="Shipping country" value={user.attributes.country} />
        <ReadOnlyField label="Payment method" value="Visa ending in 4242" />
      </div>
      <PrimaryButton disabled={conversionStatus === "recording"} onClick={onCheckoutComplete}>
        {conversionStatus === "recording" ? "Recording order" : `Place order ${formatCurrency(total)}`}
      </PrimaryButton>
      <ConversionStatus status={conversionStatus} />
    </section>
  )
}

function NewCheckout({ user, shipping, total, conversionStatus, onCheckoutComplete }) {
  return (
    <section className="rounded-lg border border-blue-200 bg-white p-6">
      <SectionHeading eyebrow="Flagged experience" title="One-page checkout" icon={ShieldCheck} />
      <div className="mb-5 grid gap-2 sm:grid-cols-3">
        {["Apple Pay", "Shop Pay", "Card"].map((method) => (
          <button
            key={method}
            className="inline-flex min-h-10 items-center justify-center gap-2 rounded-md border border-blue-200 bg-blue-50 px-4 text-sm font-semibold text-blue-700 transition hover:bg-blue-100"
          >
            <CreditCard className="h-4 w-4" />
            {method}
          </button>
        ))}
      </div>
      <div className="grid gap-3">
        <ReviewRow icon={UserRound} label="Customer" value={user.name} />
        <ReviewRow icon={Truck} label="Delivery" value={shipping === 0 ? "Free priority" : "Standard"} />
        <ReviewRow icon={CreditCard} label="Due today" value={formatCurrency(total)} />
      </div>
      <PrimaryButton disabled={conversionStatus === "recording"} onClick={onCheckoutComplete}>
        {conversionStatus === "recording" ? "Recording checkout" : "Complete secure checkout"}
      </PrimaryButton>
      <ConversionStatus status={conversionStatus} />
    </section>
  )
}

function SectionHeading({ eyebrow, title, icon: Icon }) {
  return (
    <div className="mb-6">
      <p className="flex items-center gap-2 text-xs font-semibold uppercase text-slate-500">
        {Icon && <Icon className="h-3.5 w-3.5" />}
        {eyebrow}
      </p>
      <h2 className="mt-1 text-2xl font-semibold tracking-normal text-slate-950">{title}</h2>
    </div>
  )
}

function PrimaryButton({ children, disabled, onClick }) {
  return (
    <button
      className="mt-6 inline-flex min-h-11 w-full items-center justify-center gap-2 rounded-md bg-blue-600 px-5 text-sm font-semibold text-white transition hover:bg-blue-700 disabled:cursor-not-allowed disabled:bg-slate-300"
      disabled={disabled}
      onClick={onClick}
    >
      <CheckCircle2 className="h-4 w-4" />
      {children}
    </button>
  )
}

function ConversionStatus({ status }) {
  if (status === "idle") return null

  const styles = {
    recording: "border-slate-200 bg-slate-50 text-slate-600",
    recorded: "border-blue-200 bg-blue-50 text-blue-700",
    failed: "border-red-200 bg-red-50 text-red-700"
  }

  const messages = {
    recording: "Recording conversion event",
    recorded: "Conversion event recorded",
    failed: "Conversion event failed"
  }

  return (
    <p className={`mt-4 inline-flex items-center gap-2 rounded-md border px-3 py-2 text-sm font-medium ${styles[status]}`}>
      <CheckCircle2 className="h-4 w-4" />
      {messages[status]}
    </p>
  )
}

function ReadOnlyField({ label, value }) {
  return (
    <label className="grid gap-2 text-sm font-medium text-slate-600">
      {label}
      <input
        className="min-h-10 w-full rounded-md border border-slate-200 bg-slate-50 px-3 text-sm font-medium text-slate-950"
        value={value}
        readOnly
      />
    </label>
  )
}

function ReviewRow({ label, value, icon: Icon }) {
  return (
    <div className="flex min-h-12 items-center justify-between gap-4 rounded-md border border-slate-200 bg-slate-50 px-4 py-3">
      <span className="flex items-center gap-2 text-sm text-slate-600">
        {Icon && <Icon className="h-4 w-4 text-blue-600" />}
        {label}
      </span>
      <strong className="text-right text-sm font-semibold text-slate-950">{value}</strong>
    </div>
  )
}

function OrderSummary({ subtotal, shipping, total }) {
  return (
    <aside className="rounded-lg border border-slate-200 bg-white p-6">
      <h2 className="flex items-center gap-2 text-xl font-semibold tracking-normal text-slate-950">
        <Package className="h-5 w-5 text-blue-600" />
        Order summary
      </h2>
      <div className="my-5 grid gap-4">
        {cartItems.map((item) => (
          <div className="grid grid-cols-[48px_1fr_auto] items-center gap-3" key={item.name}>
            <div className={`aspect-square w-12 rounded-md bg-gradient-to-br ${item.swatch}`} aria-hidden="true" />
            <div>
              <strong className="block text-sm font-semibold text-slate-950">{item.name}</strong>
              <span className="text-sm text-slate-500">{item.color}</span>
            </div>
            <p className="m-0 text-sm font-semibold text-slate-950">{formatCurrency(item.price)}</p>
          </div>
        ))}
      </div>
      <div className="grid gap-3 border-t border-slate-200 pt-5">
        <TotalRow label="Subtotal" value={formatCurrency(subtotal)} />
        <TotalRow label="Shipping" value={shipping === 0 ? "Free" : formatCurrency(shipping)} />
        <TotalRow label="Total" value={formatCurrency(total)} large />
      </div>
    </aside>
  )
}

function TotalRow({ label, value, large = false }) {
  return (
    <div className={`flex items-center justify-between gap-4 ${large ? "text-lg" : "text-sm"}`}>
      <span className="text-slate-600">{label}</span>
      <strong className="font-semibold text-slate-950">{value}</strong>
    </div>
  )
}

createRoot(document.getElementById("root")).render(<App />)
