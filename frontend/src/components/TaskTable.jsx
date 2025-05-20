export default function TaskTable({ tasks, onRetry }) {
  return (
    <div className="overflow-x-auto mt-8">
      <div className="min-w-[800px] w-fit mx-auto">
        <table className="table-fixed w-full bg-white shadow rounded-xl border border-gray-200">
          <thead className="bg-blue-600 text-white text-sm uppercase tracking-wider">
            <tr>
              <th className="px-4 py-3 text-left w-[80px]">ID</th>
              <th className="px-4 py-3 text-left w-[250px]">Data</th>
              <th className="px-4 py-3 text-left w-[100px]">Status</th>
              <th className="px-4 py-3 text-left w-[200px]">Updated</th>
              <th className="px-4 py-3 text-center w-[100px]">Action</th>
            </tr>
          </thead>
          <tbody>
            {tasks.length === 0 ? (
              <tr>
                <td colSpan={5} className="px-6 py-6 text-center text-gray-500 italic">
                  No tasks found.
                </td>
              </tr>
            ) : (
              tasks.map((task) => (
                <tr
                  key={task.id}
                  className="border-b last:border-none hover:bg-gray-50 transition"
                >
                  <td className="px-4 py-3 text-sm font-mono text-gray-700 truncate">{task.id}</td>
                  <td className="px-4 py-3 text-sm text-gray-800 break-words">{task.data}</td>
                  <td
                    className={`px-4 py-3 text-sm font-semibold ${
                      task.status === 'done'
                        ? 'text-green-600'
                        : task.status === 'failed'
                        ? 'text-red-600'
                        : 'text-yellow-600'
                    }`}
                  >
                    {task.status}
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-500">
                    {new Date(task.updated_at).toLocaleString()}
                  </td>
                  <td className="px-4 py-3 text-center">
                    {task.status === 'failed' && (
                      <button
                        onClick={() => onRetry(task.id)}
                        className="px-4 py-1.5 bg-red-500 hover:bg-red-600 text-white text-xs font-medium rounded-md transition"
                        type="button"
                      >
                        Retry
                      </button>
                    )}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
