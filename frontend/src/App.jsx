import React from "react";
import TaskForm from "./components/TaskForm";
import TaskList from "./components/TaskList";

function App() {
  return (
    <div style={{ padding: "2rem" }}>
      <h1>Task Queue Dashboard</h1>
      <TaskForm />
      <TaskList />
    </div>
  );
}

export default App;