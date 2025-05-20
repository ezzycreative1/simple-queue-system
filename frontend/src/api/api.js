const BASE_URL = 'http://localhost:8081/api';

export async function fetchTasks(status = '', page = 1, limit = 10) {
  const params = new URLSearchParams({ status, page, limit });
  const res = await fetch(`${BASE_URL}/queue?${params.toString()}`);
  return res.json();
}

export async function retryTask(id) {
  const res = await fetch(`${BASE_URL}/retry/${id}`, {
    method: 'POST',
  });
  return res.json();
}
