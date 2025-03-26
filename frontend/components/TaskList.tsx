import { useState } from "react";

interface Task {
  id: string;
  name: string;
  description: string;
  dueDate: string;
  assignedto: string;
  userEmail: string;
  status: boolean;
}

interface TaskListProps {
  tasks: Task[];
  onStatusChange: (updatedTask: Task) => void;
}

export default function TaskList({ tasks, onStatusChange }: TaskListProps) {
  const [error, setError] = useState<string | null>(null);

  const toggleStatus = async (taskId: string) => {
    try {
      const response = await fetch(`http://localhost:8080/api/tasks/${taskId}/status`, {
        method: 'PATCH',
        credentials: 'include',
      });
      if (!response.ok) throw new Error('Failed to update task status');
      const data = await response.json();
      onStatusChange({ ...tasks.find((task) => task.id === taskId)!, status: data.status });
    } catch {
      setError('Failed to update task status');
    }
  };

  return (
    <div className="bg-white p-8 rounded-lg shadow-lg overflow-hidden">
      {error && <p className="text-red-500">{error}</p>}
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">
                Task No.
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Task Name
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Description
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Due Date
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Assigned To
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status
                <div className="inline-block relative group ml-2">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="h-4 w-4 text-gray-500 cursor-pointer"
                    viewBox="0 0 20 20"
                    fill="currentColor"
                  >
                    <path d="M10 18a8 8 0 100-16 8 8 0 000 16zm0-14a6 6 0 110 12 6 6 0 010-12zm.75 9.75a.75.75 0 11-1.5 0v-1.5a.75.75 0 011.5 0v1.5zm0-4.5a.75.75 0 11-1.5 0v-3a.75.75 0 011.5 0v3z" />
                  </svg>
                  <span className="absolute -top-1 left-1/2 transform -translate-x-1/2 bg-gray-800/80 text-white text-xs px-2 py-1 rounded opacity-0 group-hover:opacity-100 transition-opacity z-500">
                    Click to update status
                  </span>
                </div>
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {tasks.length === 0 ? (
              <tr>
                <td colSpan={6} className="px-6 py-4 text-center">
                  No tasks found
                </td>
              </tr>
            ) : (
              tasks.map((task, index) => (
                <tr key={task.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 text-center">{index + 1}</td>
                  <td className="px-6 py-4">{task.name}</td>
                  <td className="px-6 py-4">{task.description}</td>
                  <td className="px-6 py-4">{new Date(task.dueDate).toLocaleDateString()}</td>
                  <td className="px-6 py-4">{task.assignedto}</td>
                  <td className="px-6 py-4">
                    <button
                      className={`px-4 py-2 rounded ${
                        task.status ? "bg-green-500 text-white" : "bg-gray-500 text-white"
                      }`}
                      onClick={() => toggleStatus(task.id)}
                    >
                      {task.status ? "Completed" : "Incomplete"}
                    </button>
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
