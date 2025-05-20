import React, { useState } from "react";

function TaskForm() {
  const [data, setData] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();
    const res = await fetch("/api/enqueue", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ data }),
    });

    const result = await res.json();
    alert(result.message || result.status);
    setData("");
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        value={data}
        onChange={(e) => setData(e.target.value)}
        placeholder="Enter task data"
        required
      />
      <button type="submit">Enqueue</button>
    </form>
  );
}

export default TaskForm;
