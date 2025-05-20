import { useEffect, useState, useRef } from 'react';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { fetchTasks, retryTask } from './api/api';
import TaskTable from './components/TaskTable';
import StatusFilter from './components/StatusFilter';
import Pagination from './components/Pagination';

export default function App() {
  const [tasks, setTasks] = useState([]);
  const [status, setStatus] = useState('');
  const [page, setPage] = useState(1);
  const [limit] = useState(10);
  const [total, setTotal] = useState(0);
  const dataInputRef = useRef();

  const loadTasks = async () => {
    try {
      const res = await fetchTasks(status, page, limit);
      setTasks(res.data.tasks);
      setTotal(res.data.meta.total);
    } catch (err) {
      toast.error("Gagal memuat data");
    }
  };

  const handleRetry = async (id) => {
    const res = await retryTask(id);
    if (res.status === 'success') {
      toast.success("Task berhasil di-retry");
      loadTasks();
    } else {
      toast.error(res.message || "Retry gagal");
    }
  };

  const handleAddTask = async (e) => {
    e.preventDefault();
    const data = dataInputRef.current.value.trim();
    if (!data) {
      toast.error("Data task harus diisi");
      return;
    }

    try {
      const res = await fetch('/api/enqueue', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ data }),
      });
      const json = await res.json();
      if (json.status === 'success') {
        toast.success("Task berhasil ditambahkan");
        dataInputRef.current.value = '';
        loadTasks(); // refresh list
      } else {
        toast.error(json.message || "Gagal tambah task");
      }
    } catch (err) {
      toast.error("Error saat menambahkan task");
    }
  };

  useEffect(() => {
    loadTasks();
    const interval = setInterval(loadTasks, 5000); // Auto-refresh setiap 5 detik
    return () => clearInterval(interval);
  }, [status, page]);

  return (
    <div className="min-h-screen bg-gradient-to-br from-white to-blue-50 text-gray-900 p-6 max-w-5xl mx-auto font-sans">
      <h1 className="text-4xl font-bold mb-6 text-center tracking-tight text-blue-700 drop-shadow-sm">
        Task Queue
      </h1>

      <div className="flex flex-col md:flex-row md:justify-between md:items-center gap-4 mb-4">
        <StatusFilter status={status} setStatus={setStatus} />
      </div>

      <form onSubmit={handleAddTask} className="flex flex-col md:flex-row gap-3 mb-6">
        <input
          ref={dataInputRef}
          type="text"
          placeholder="Masukkan data task"
          className="flex-1 p-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-400 text-sm placeholder-gray-400"
        />
        <button
          type="submit"
          className="px-5 py-2.5 bg-blue-600 hover:bg-blue-700 text-black text-sm font-medium rounded-lg shadow-md transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-blue-400"
        >
          ➕ Add Task
        </button>
      </form>

      <TaskTable tasks={tasks} onRetry={handleRetry} />
      <Pagination page={page} total={total} limit={limit} onPageChange={setPage} />

      <ToastContainer />
    </div>
  );
}
