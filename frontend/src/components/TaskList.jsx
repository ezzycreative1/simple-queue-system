import React, { useEffect, useState } from "react";

export default function TaskQueueApp() {
  const [tasks, setTasks] = useState([]);
  const [data, setData] = useState("");
  const [filter, setFilter] = useState("");
  const [page, setPage] = useState(1);
  const pageSize = 10;

  const fetchTasks = async () => {
    const res = await fetch("http://localhost:8081/api/queue");
    const json = await res.json();
    setTasks(json.data);
  };

  useEffect(() => {
    fetchTasks();
    const interval = setInterval(fetchTasks, 3000);
    return () => clearInterval(interval);
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!data.trim()) return;
    await fetch("http://localhost:8081/api/enqueue", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ data })
    });
    setData("");
    fetchTasks();
  };

  const retryTask = async (id) => {
    await fetch(`http://localhost:8081/api/retry/${id}`, {
      method: "POST"
    });
    fetchTasks();
  };

  const filteredTasks = tasks.filter(t => !filter || t.status === filter);
  const pagedTasks = filteredTasks.slice((page - 1) * pageSize, page * pageSize);

  return (
    <div className="p-4 max-w-4xl mx-auto">
      <div className="mb-4">
        <label className="mr-2 font-medium">Filter Status:</label>
        <select value={filter} onChange={e => { setFilter(e.target.value); setPage(1); }} className="p-2 border rounded">
          <option value="">Semua</option>
          <option value="pending">Pending</option>
          <option value="processing">Processing</option>
          <option value="done">Done</option>
          <option value="failed">Failed</option>
        </select>
      </div>

      <table className="w-full border text-sm">
        <thead>
          <tr className="bg-gray-100 text-left">
            <th className="p-2">ID</th>
            <th className="p-2">Data</th>
            <th className="p-2">Status</th>
            <th className="p-2">Created At</th>
            <th className="p-2">Aksi</th>
          </tr>
        </thead>
        <tbody>
          {pagedTasks.map(task => (
            <tr key={task.id} className="border-t">
              <td className="p-2">{task.id}</td>
              <td className="p-2">{task.data}</td>
              <td className="p-2 capitalize">{task.status}</td>
              <td className="p-2">{new Date(task.created_at).toLocaleString()}</td>
              <td className="p-2">
                {task.status === "failed" && (
                  <button
                    onClick={() => retryTask(task.id)}
                    className="bg-yellow-500 text-white px-3 py-1 rounded"
                  >
                    Retry
                  </button>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      <div className="flex justify-between mt-4">
        <button
          className="px-4 py-2 bg-gray-200 rounded disabled:opacity-50"
          onClick={() => setPage(p => p - 1)}
          disabled={page === 1}
        >
          Prev
        </button>

        <span className="text-sm mt-2">Halaman {page}</span>

        <button
          className="px-4 py-2 bg-gray-200 rounded disabled:opacity-50"
          onClick={() => setPage(p => p + 1)}
          disabled={page * pageSize >= filteredTasks.length}
        >
          Next
        </button>
      </div>
    </div>
  );
}
