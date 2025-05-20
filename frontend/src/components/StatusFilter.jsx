export default function StatusFilter({ status, setStatus }) {
  return (
    <select value={status} onChange={(e) => setStatus(e.target.value)} className="p-2 rounded border">
      <option value="">All</option>
      <option value="pending">Pending</option>
      <option value="processing">Processing</option>
      <option value="done">Done</option>
      <option value="failed">Failed</option>
    </select>
  );
}