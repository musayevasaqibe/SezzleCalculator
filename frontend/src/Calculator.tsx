import React, { useState } from "react";

type Resp = { result?: number | null; error?: string | null };

export default function Calculator() {
  const [display, setDisplay] = useState<string>(""); // current typed number (right-hand)
  const [stored, setStored] = useState<number | null>(null); // left-hand operand
  const [pendingOp, setPendingOp] = useState<string | null>(null); // add, sub, etc.
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const digits = ["7", "8", "9", "4", "5", "6", "1", "2", "3", "0"];
  const ops = [
    { value: "add", label: "+" },
    { value: "sub", label: "-" },
    { value: "mul", label: "×" },
    { value: "div", label: "÷" },
  ];

  const opSymbols: Record<string, string> = {
    add: "+",
    sub: "-",
    mul: "×",
    div: "÷",
    pow: "^",
    pct: "%",
  };

  function pushDigit(d: string) {
    setDisplay((prev) => (prev === "0" ? d : prev + d));
    setError(null);
  }

  function clearAll() {
    setDisplay("");
    setStored(null);
    setPendingOp(null);
    setError(null);
  }

  function deleteLast() {
    setError(null);
    if (display !== "") {
      const next = display.slice(0, -1);
      setDisplay(next);
      return;
    }
    // If display is empty but there's a stored value and no pendingOp,
    // allow deleting last digit from stored (optional UX).
    if (stored !== null && !pendingOp) {
      const s = String(stored);
      if (s.length <= 1) {
        setStored(null);
      } else {
        // preserve numeric value after removing last char
        const trimmed = s.slice(0, -1);
        const num = Number(trimmed);
        if (!Number.isNaN(num)) {
          setStored(num);
        } else {
          setStored(null);
        }
      }
    }
  }

  function chooseOp(op: string) {
    if (display !== "") {
      setStored(Number(display));
      setDisplay("");
      setPendingOp(op);
    } else if (stored !== null) {
      setPendingOp(op);
    } else {
      setStored(0);
      setPendingOp(op);
    }
    setError(null);
  }

  async function computeRemote(body: any) {
    setLoading(true);
    setError(null);
    try {
      const res = await fetch("http://localhost:8080/api/calc", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body),
      });
      const data: Resp = await res.json();
      if (data.error) {
        setError(data.error);
      } else {
        setDisplay(String(data.result ?? ""));
        setStored(null);
        setPendingOp(null);
      }
    } catch (err: any) {
      setError(err?.message ?? "Network error");
    } finally {
      setLoading(false);
    }
  }

  function pressEquals() {
    if (!pendingOp) {
      setError("No operation selected");
      return;
    }
    if (stored === null || display === "") {
      setError("Missing operand");
      return;
    }
    computeRemote({ op: pendingOp, a: stored, b: Number(display) });
  }

  // Build the full display text: stored + operator + current display
  const displayText = (() => {
    if (stored !== null) {
      const left = String(stored);
      const op = pendingOp ? opSymbols[pendingOp] || pendingOp : "";
      const right = display || "";
      return left + (op ? op : "") + right;
    }
    return display || "0";
  })();

  return (
    <div className="calc-wrapper">
      <div className="display" data-testid="display">
        {displayText}
      </div>

      <div className="button-row">
        <button className="btn" onClick={deleteLast} disabled={loading}>
          DEL
        </button>
        <button className="btn wide" onClick={clearAll} disabled={loading}>
          C
        </button>
        <button
          className={`btn op ${pendingOp === "div" ? "active" : ""}`}
          onClick={() => chooseOp("div")}
          disabled={loading}
        >
          ÷
        </button>
        <button
          className={`btn op ${pendingOp === "mul" ? "active" : ""}`}
          onClick={() => chooseOp("mul")}
          disabled={loading}
        >
          ×
        </button>
      </div>

      <div className="grid">
        <div className="digits">
          {digits.map((d) => (
            <button
              key={d}
              className="btn digit"
              onClick={() => pushDigit(d)}
              disabled={loading}
            >
              {d}
            </button>
          ))}
        </div>

        <div className="ops-col">
          {ops.map((o) => (
            <button
              key={o.value}
              className={`btn op ${pendingOp === o.value ? "active" : ""}`}
              onClick={() => chooseOp(o.value)}
              disabled={loading}
            >
              {o.label}
            </button>
          ))}
          <button className="btn eq" onClick={pressEquals} disabled={loading}>
            =
          </button>
        </div>
      </div>

      <div className="status">
        {loading && <span>Calculating...</span>}
        {error && <span className="error">Error: {error}</span>}
      </div>
    </div>
  );
}
