import React, { useState } from "react";

type Resp = { result?: number | null; error?: string | null };

export default function Calculator() {
  const [display, setDisplay] = useState<string>(""); // what's shown on screen
  const [stored, setStored] = useState<number | null>(null); // stored operand A
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

  return (
    <div className="calc-wrapper">
      <div className="display" data-testid="display">
        {display || (stored !== null ? String(stored) : "0")}
      </div>

      <div className="button-row">
        <button className="btn wide" onClick={clearAll} disabled={loading}>
          C
        </button>
        <button
          className="btn op"
          onClick={() => chooseOp("div")}
          disabled={loading}
        >
          ÷
        </button>
        <button
          className="btn op"
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
              className="btn op"
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
