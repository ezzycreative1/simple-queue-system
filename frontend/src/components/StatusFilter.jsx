export default function StatusFilter({ status, setStatus }) {
  return (
    <select
      value={status}
      onChange={(e) => setStatus(e.target.value)}
      className="px-3 py-2 border rounded text-gray-700 dark:text-white dark:bg-gray-800 dark:border-gray-600"
    >
      <option value="">All</option>
      <option value="waiting">Waiting</option>
      <option value="done">Done</option>
      <option value="failed">Failed</option>
    </select>
  );
}
