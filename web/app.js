const $ = (id) => document.getElementById(id);
const money = new Intl.NumberFormat("en-US", { style: "currency", currency: "USD", maximumFractionDigits: 2 });
const fields = ["homePrice", "downPercent", "annualRate", "termYears", "taxPercent", "annualIns"];
const labels = { monthly: "monthly", quarterly: "quarterly", annual: "annual" };

let period = "monthly";
let lastResult = null;   
let timer = null;        
let controller = null;  

function readForm() {
  const num = (id) => (document.getElementById(id).value === "" ? NaN : Number(document.getElementById(id).value));
  return {
    homePrice: num("homePrice"),
    downPercent: num("downPercent"),
    annualRate: num("annualRate"),
    termYears: Math.trunc(num("termYears")),
    taxPercent: num("taxPercent"),
    annualIns: num("annualIns"),
  };
}

async function calculate() {
  const body = readForm();
  if (Object.values(body).some(Number.isNaN)) {
    showError("Fill in every field with a number.");
    return;
  }
  if (controller) controller.abort();
  controller = new AbortController();

  try {
    const res = await fetch("/api/calculate", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
      signal: controller.signal,
    });
    const data = await res.json();
    if (!res.ok) {
      showError(data.error || "Something went wrong.");
      return;
    }
    lastResult = data;
    showError("");
    render();
  } catch (err) {
    if (err.name === "AbortError") return; // a newer request replaced this one
    showError("Cannot reach the server. Is it running?");
  }
}

function showError(msg) {
  $("error").textContent = msg;
  if (msg) {
    $("results").hidden = true;
    $("loanCard").hidden = true;
  }
}

function render() {
  const r = lastResult;
  if (!r) return;
  const p = r[period];
  $("results").hidden = false;
  $("loanCard").hidden = false;
  $("periodLabel").textContent = "Estimated " + labels[period] + " payment";
  $("total").textContent = money.format(p.total);
  $("pi").textContent = money.format(p.principalAndInterest);
  $("tax").textContent = money.format(p.propertyTax);
  $("ins").textContent = money.format(p.insurance);
  const share = p.total > 0 ? (p.principalAndInterest / p.total) * 100 : 100;
  $("b1").style.width = share + "%";
  $("b2").style.width = 100 - share + "%";

  $("gM").textContent = money.format(r.monthly.total);
  $("gQ").textContent = money.format(r.quarterly.total);
  $("gA").textContent = money.format(r.annual.total);
  $("down").textContent = money.format(r.downPayment);
  $("loan").textContent = money.format(r.loanAmount);
  $("interest").textContent = money.format(r.totalInterest);
  $("repay").textContent = money.format(r.totalRepayment);
}

fields.forEach((id) =>
  $(id).addEventListener("input", () => {
    clearTimeout(timer);
    timer = setTimeout(calculate, 250);
  })
);

document.querySelectorAll(".tabs button").forEach((btn) =>
  btn.addEventListener("click", () => {
    period = btn.dataset.period;
    document.querySelectorAll(".tabs button").forEach((b) => b.setAttribute("aria-pressed", String(b === btn)));
    render();
  })
);

calculate();