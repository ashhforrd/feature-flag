import React, { useEffect, useMemo, useState } from "react"
import { createRoot } from "react-dom/client"
import CountUpPackage from "react-countup"
import {
  BarChart3,
  CheckCircle2,
  FlaskConical,
  Gauge,
  GitBranch,
  ListFilter,
  Percent,
  Power,
  RefreshCw,
  Save,
  Settings2,
  ShieldCheck,
  SlidersHorizontal,
  TrendingUp
} from "lucide-react"
import "./styles.css"

const CountUp = CountUpPackage.default || CountUpPackage
const API_BASE_URL = import.meta.env.VITE_CONFIG_SERVICE_URL || "/api"

async function api(path, options = {}) {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...options.headers
    },
    ...options
  })

  if (!response.ok) throw new Error(`Request failed with ${response.status}`)
  return response.json()
}

function formatPercent(value) {
  return `${Math.round(value * 1000) / 10}%`
}

function App() {
  const [flags, setFlags] = useState([])
  const [selectedKey, setSelectedKey] = useState("")
  const [selectedFlag, setSelectedFlag] = useState(null)
  const [exposures, setExposures] = useState(null)
  const [results, setResults] = useState(null)
  const [draft, setDraft] = useState({ enabled: false, rolloutPercentage: 0 })
  const [status, setStatus] = useState("loading")
  const [saveStatus, setSaveStatus] = useState("idle")

  async function loadFlags() {
    setStatus("loading")
    try {
      const data = await api("/flags")
      setFlags(data)
      setSelectedKey((current) => current || data[0]?.key || "")
      setStatus("ready")
    } catch {
      setStatus("error")
    }
  }

  async function loadSelectedFlag(key) {
    if (!key) {
      setSelectedFlag(null)
      setExposures(null)
      setResults(null)
      return
    }

    setStatus("loading-detail")
    try {
      const [flag, exposureSummary, experimentResult] = await Promise.all([
        api(`/flags/${key}`),
        api(`/flags/${key}/exposures`),
        api(`/flags/${key}/results`)
      ])

      setSelectedFlag(flag)
      setDraft({ enabled: flag.enabled, rolloutPercentage: flag.rolloutPercentage })
      setExposures(exposureSummary)
      setResults(experimentResult)
      setStatus("ready")
    } catch {
      setStatus("error")
    }
  }

  useEffect(() => {
    loadFlags()
  }, [])

  useEffect(() => {
    loadSelectedFlag(selectedKey)
  }, [selectedKey])

  async function saveFlag() {
    if (!selectedFlag) return

    setSaveStatus("saving")
    try {
      const updated = await api(`/flags/${selectedFlag.key}`, {
        method: "PATCH",
        body: JSON.stringify({
          enabled: draft.enabled,
          rolloutPercentage: Number(draft.rolloutPercentage)
        })
      })

      setSelectedFlag(updated)
      setFlags((current) => current.map((flag) => (flag.key === updated.key ? updated : flag)))
      setSaveStatus("saved")
      window.setTimeout(() => setSaveStatus("idle"), 1600)
    } catch {
      setSaveStatus("error")
    }
  }

  const selectedIndex = useMemo(
    () => flags.findIndex((flag) => flag.key === selectedKey),
    [flags, selectedKey]
  )

  return (
    <main className="min-h-screen bg-slate-50 text-slate-950">
      <div className="mx-auto grid min-h-screen max-w-7xl grid-cols-1 lg:grid-cols-[280px_1fr]">
        <aside className="border-b border-slate-200 bg-white p-5 lg:border-b-0 lg:border-r">
          <div className="mb-6">
            <p className="flex items-center gap-2 text-xs font-semibold uppercase text-slate-500">
              <ShieldCheck className="h-3.5 w-3.5" />
              Feature Platform
            </p>
            <h1 className="mt-1 text-2xl font-semibold tracking-normal text-slate-950">Dashboard</h1>
          </div>

          <button
            className="mb-4 inline-flex min-h-10 w-full items-center justify-center gap-2 rounded-md border border-slate-200 bg-white px-4 text-sm font-medium text-slate-700 transition hover:border-slate-400 hover:bg-slate-50"
            onClick={loadFlags}
          >
            <RefreshCw className="h-4 w-4" />
            Refresh flags
          </button>

          <nav className="grid gap-2" aria-label="Flags">
            {flags.map((flag, index) => (
              <button
                key={flag.key}
                className={`rounded-md border p-3 text-left transition ${
                  flag.key === selectedKey
                    ? "border-slate-900 bg-slate-900 text-white"
                    : "border-slate-200 bg-white text-slate-950 hover:border-slate-400"
                }`}
                onClick={() => setSelectedKey(flag.key)}
              >
                <div className="mb-2 flex items-center justify-between gap-3">
                  <strong className="flex min-w-0 items-center gap-2 truncate text-sm font-semibold">
                    <GitBranch className="h-4 w-4 shrink-0" />
                    <span className="truncate">{flag.key}</span>
                  </strong>
                  <span className="text-xs opacity-70">#{index + 1}</span>
                </div>
                <div className="flex items-center justify-between gap-3 text-xs opacity-80">
                  <span>{flag.enabled ? "Enabled" : "Disabled"}</span>
                  <span>{flag.rolloutPercentage}% rollout</span>
                </div>
              </button>
            ))}
          </nav>
        </aside>

        <section className="p-5 sm:p-6">
          <header className="mb-5 flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
            <div>
              <p className="flex items-center gap-2 text-xs font-semibold uppercase text-slate-500">
                <ListFilter className="h-3.5 w-3.5" />
                Selected flag
              </p>
              <h2 className="mt-1 text-3xl font-semibold tracking-normal text-slate-950">
                {selectedFlag?.name || selectedKey || "No flag selected"}
              </h2>
              {selectedFlag && <p className="mt-1 max-w-2xl text-sm text-slate-600">{selectedFlag.description}</p>}
            </div>
            <div className="inline-flex w-fit items-center gap-2 rounded-full border border-slate-200 bg-white px-3 py-1.5 text-sm font-medium text-slate-600">
              <Gauge className="h-4 w-4" />
              {status === "loading" || status === "loading-detail" ? "Loading" : status === "error" ? "API error" : "Live"}
            </div>
          </header>

          {selectedFlag ? (
            <div className="grid gap-4">
              <section className="grid gap-4 xl:grid-cols-[360px_1fr]">
                <FlagControls draft={draft} setDraft={setDraft} saveFlag={saveFlag} saveStatus={saveStatus} />
                <FlagSnapshot flag={selectedFlag} selectedIndex={selectedIndex} />
              </section>

              <section className="grid gap-4 md:grid-cols-3">
                <MetricCard icon={BarChart3} label="Total exposures" value={exposures?.total || 0} />
                <MetricCard icon={TrendingUp} label="Enabled exposures" value={exposures?.enabled || 0} />
                <MetricCard icon={Gauge} label="Disabled exposures" value={exposures?.disabled || 0} />
              </section>

              <section className="grid gap-4 xl:grid-cols-2">
                <VariantCard title="Enabled variant" variant={results?.enabled} />
                <VariantCard title="Disabled variant" variant={results?.disabled} />
              </section>
            </div>
          ) : (
            <div className="rounded-lg border border-slate-200 bg-white p-6 text-sm text-slate-600">
              Create a flag in the config service to see dashboard data.
            </div>
          )}
        </section>
      </div>
    </main>
  )
}

function FlagControls({ draft, setDraft, saveFlag, saveStatus }) {
  return (
    <section className="rounded-lg border border-slate-200 bg-white p-5">
      <SectionHeading icon={SlidersHorizontal} eyebrow="Controls" title="Release settings" />

      <label className="mb-4 flex items-center justify-between gap-4 rounded-md border border-slate-200 bg-slate-50 p-3">
        <div>
          <span className="flex items-center gap-2 text-sm font-semibold text-slate-950">
            <Power className="h-4 w-4" />
            Enabled
          </span>
          <span className="text-sm text-slate-500">Kill switch for this flag</span>
        </div>
        <input
          className="h-4 w-4 accent-slate-900"
          type="checkbox"
          checked={draft.enabled}
          onChange={(event) => setDraft((current) => ({ ...current, enabled: event.target.checked }))}
        />
      </label>

      <label className="grid gap-3">
        <div className="flex items-center justify-between gap-4">
          <span className="flex items-center gap-2 text-sm font-semibold text-slate-950">
            <Percent className="h-4 w-4" />
            Rollout percentage
          </span>
          <strong className="text-sm font-semibold text-slate-950">{draft.rolloutPercentage}%</strong>
        </div>
        <input
          className="accent-slate-900"
          type="range"
          min="0"
          max="100"
          value={draft.rolloutPercentage}
          onChange={(event) => setDraft((current) => ({ ...current, rolloutPercentage: Number(event.target.value) }))}
        />
      </label>

      <button
        className="mt-5 inline-flex min-h-10 w-full items-center justify-center gap-2 rounded-md bg-slate-900 px-4 text-sm font-semibold text-white transition hover:bg-slate-700 disabled:cursor-not-allowed disabled:bg-slate-300"
        disabled={saveStatus === "saving"}
        onClick={saveFlag}
      >
        <Save className="h-4 w-4" />
        {saveStatus === "saving" ? "Saving" : saveStatus === "saved" ? "Saved" : "Save changes"}
      </button>
      {saveStatus === "error" && <p className="mt-3 text-sm font-medium text-red-700">Failed to save changes.</p>}
    </section>
  )
}

function FlagSnapshot({ flag, selectedIndex }) {
  return (
    <section className="rounded-lg border border-slate-200 bg-white p-5">
      <div className="mb-5 flex items-start justify-between gap-4">
        <SectionHeading icon={Settings2} eyebrow="Configuration" title="Flag snapshot" />
        <span className="rounded-full bg-slate-100 px-3 py-1 text-xs font-medium text-slate-600">#{selectedIndex + 1}</span>
      </div>
      <dl className="grid gap-3 sm:grid-cols-2">
        <InfoItem label="Key" value={flag.key} />
        <InfoItem label="State" value={flag.enabled ? "Enabled" : "Disabled"} />
        <InfoItem label="Rollout" value={`${flag.rolloutPercentage}%`} />
        <InfoItem label="Rules" value={`${flag.targetingRules?.length || 0}`} />
      </dl>
    </section>
  )
}

function SectionHeading({ eyebrow, title, icon: Icon }) {
  return (
    <div>
      <p className="flex items-center gap-2 text-xs font-semibold uppercase text-slate-500">
        {Icon && <Icon className="h-3.5 w-3.5" />}
        {eyebrow}
      </p>
      <h3 className="mt-1 text-xl font-semibold tracking-normal text-slate-950">{title}</h3>
    </div>
  )
}

function InfoItem({ label, value }) {
  return (
    <div className="rounded-md border border-slate-200 bg-slate-50 p-3">
      <dt className="text-xs font-semibold uppercase text-slate-500">{label}</dt>
      <dd className="mt-2 break-words text-base font-semibold text-slate-950">{value}</dd>
    </div>
  )
}

function MetricCard({ label, value, icon: Icon }) {
  return (
    <article className="rounded-lg border border-slate-200 bg-white p-5">
      <p className="mb-2 flex items-center gap-2 text-sm font-medium text-slate-500">
        {Icon && <Icon className="h-4 w-4" />}
        {label}
      </p>
      <strong className="text-4xl font-semibold tracking-normal text-slate-950">
        <CountUp end={value} duration={0.7} preserveValue />
      </strong>
    </article>
  )
}

function VariantCard({ title, variant }) {
  const data = variant || { exposures: 0, conversions: 0, conversionRate: 0 }
  return (
    <article className="rounded-lg border border-slate-200 bg-white p-5">
      <div className="mb-5 flex items-start justify-between gap-4">
        <SectionHeading icon={FlaskConical} eyebrow="Experiment result" title={title} />
        <span className="inline-flex items-center gap-1.5 rounded-full border border-slate-200 bg-slate-50 px-3 py-1 text-xs font-semibold text-slate-700">
          <TrendingUp className="h-3.5 w-3.5" />
          {formatPercent(data.conversionRate)}
        </span>
      </div>
      <div className="grid gap-3 sm:grid-cols-3">
        <MiniMetric label="Exposures" value={data.exposures} />
        <MiniMetric label="Conversions" value={data.conversions} />
        <MiniMetric label="Rate" value={data.conversionRate * 100} suffix="%" decimals={1} />
      </div>
    </article>
  )
}

function MiniMetric({ label, value, suffix = "", decimals = 0 }) {
  return (
    <div className="rounded-md border border-slate-200 bg-slate-50 p-3">
      <p className="mb-2 text-xs font-semibold uppercase text-slate-500">{label}</p>
      <strong className="text-xl font-semibold tracking-normal text-slate-950">
        <CountUp end={value} duration={0.7} decimals={decimals} preserveValue />{suffix}
      </strong>
    </div>
  )
}

createRoot(document.getElementById("root")).render(<App />)
