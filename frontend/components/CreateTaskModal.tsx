import { useState } from 'react';
import { useAuth } from '../app/contexts/AuthContext';

interface User {
  email: string;
  username: string;
}

interface CreateTaskModalProps {
  users: User[];
  onClose: () => void;
  onTaskCreated: () => void;
}

export default function CreateTaskModal({ users, onClose, onTaskCreated }: CreateTaskModalProps) {
  const { user } = useAuth(); // Replace useContext with useAuth hook

  const [newTask, setNewTask] = useState({
    name: '',
    description: '',
    dueDate: '',
    assignedTo: '',
    status: false,
    userEmail: user?.email || '', // Initialize with the logged-in user's email
  });
  const [errors, setErrors] = useState<{ [key: string]: string }>({});

  const handleCreateTask = async (e: React.FormEvent) => {
    e.preventDefault();
    const newErrors: any = {};
    if (!newTask.name.trim()) newErrors.name = 'Task name is required';
    if (!newTask.description.trim()) newErrors.description = 'Task description is required';
    if (!newTask.dueDate) newErrors.dueDate = 'Due date is required';
    if (!newTask.assignedTo) newErrors.assignedTo = 'Assigned user is required';
    if (!newTask.userEmail) newErrors.userEmail = 'User email is missing'; // Optional validation

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    try {
      const formattedTask = {
        ...newTask,
        dueDate: new Date(newTask.dueDate).toISOString(), // Convert dueDate to ISO 8601 format
      };

      const response = await fetch('http://localhost:8080/api/tasks', {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(formattedTask), // Use formattedTask in the payload
      });
      if (!response.ok) throw new Error('Failed to create task');
      setNewTask({ name: '', description: '', dueDate: '', assignedTo: '', status: false, userEmail: user?.email || '' });
      setErrors({});
      onClose();
      onTaskCreated();
    } catch {
      setErrors({ name: 'Failed to create task' });
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
      <div className="bg-white rounded-lg p-8 max-w-md w-full">
        <form onSubmit={handleCreateTask} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700">Task Name</label>
            <input
              type="text"
              value={newTask.name}
              onChange={(e) => setNewTask({ ...newTask, name: e.target.value })}
              className={`mt-1 block w-full rounded-md shadow-sm ${
                errors.name ? 'border-red-300 focus:border-red-500' : 'border-gray-300'
              }`}
            />
            {errors.name && <p className="mt-1 text-sm text-red-600">{errors.name}</p>}
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">Assign To</label>
            <select
              value={newTask.assignedTo}
              onChange={(e) => setNewTask({ ...newTask, assignedTo: e.target.value })}
              className={`mt-1 block w-full rounded-md shadow-sm ${
                errors.assignedTo ? 'border-red-300 focus:border-red-500' : 'border-gray-300'
              }`}
            >
              <option value="">Select a user</option>
              {users.map((user) => (
                <option key={user.email} value={user.username}>
                  {user.username}
                </option>
              ))}
            </select>
            {errors.assignedTo && <p className="mt-1 text-sm text-red-600">{errors.assignedTo}</p>}
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">Description</label>
            <textarea
              value={newTask.description}
              onChange={(e) => setNewTask({ ...newTask, description: e.target.value })}
              className={`mt-1 block w-full rounded-md shadow-sm ${
                errors.description ? 'border-red-300 focus:border-red-500' : 'border-gray-300'
              }`}
              rows={3}
            />
            {errors.description && <p className="mt-1 text-sm text-red-600">{errors.description}</p>}
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700">Due Date</label>
            <input
              type="date"
              value={newTask.dueDate}
              onChange={(e) => setNewTask({ ...newTask, dueDate: e.target.value })}
              className={`mt-1 block w-full rounded-md shadow-sm ${
                errors.dueDate ? 'border-red-300 focus:border-red-500' : 'border-gray-300'
              }`}
            />
            {errors.dueDate && <p className="mt-1 text-sm text-red-600">{errors.dueDate}</p>}
          </div>
          <div className="flex gap-4 mt-6">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded hover:bg-gray-50"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="flex-1 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            >
              Create Task
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
