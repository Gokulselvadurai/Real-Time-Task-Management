'use client';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/app/contexts/AuthContext';
import TaskList from '@/components/TaskList';
import CreateTaskModal from '@/components/CreateTaskModal';

interface Task {
  id: string;
  name: string;
  description: string;
  dueDate: string;
  assignedto: string;
  status: boolean; // Changed type from string to boolean
  userEmail: string; // Added userEmail property
}

interface User {
  email: string;
  username: string;
}

export default function Dashboard() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showModal, setShowModal] = useState(false);
  const { logout, checkAuthStatus, user } = useAuth();
  const router = useRouter();
  const [users, setUsers] = useState<User[]>([]);

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const isValid = await checkAuthStatus();
        if (!isValid) {
          router.push('/login');
          return;
        }
        await fetchTasks();
      } catch {
        setError('Authentication failed');
        router.push('/login');
      }
    };
    checkAuth();
  }, []);

  const fetchTasks = async () => {
    try {
      setIsLoading(true);
      setError(null);
      const response = await fetch('http://localhost:8080/api/tasks', {
        credentials: 'include',
      });
      if (response.status === 401) {
        await logout();
        return;
      }
      if (!response.ok) throw new Error('Failed to fetch tasks');
      const data: Task[] = await response.json(); // Ensure data matches the updated Task interface
      setTasks(data || []);
    } catch {
      setError('Failed to load tasks');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const response = await fetch('http://localhost:8080/api/users', {
          credentials: 'include',
        });
        if (response.ok) {
          const data = await response.json();
          setUsers(data);
        }
      } catch {
        console.error('Error fetching users');
      }
    };
    fetchUsers();
  }, []);

  const handleStatusChange = (updatedTask: Task) => {
    setTasks((prevTasks) =>
      prevTasks.map((task) => (task.id === updatedTask.id ? updatedTask : task))
    );
  };

  if (error) {
    return <div className="text-red-500 p-4">{error}</div>;
  }

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <div className="max-w-6xl mx-auto">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-2xl font-bold">Real Time Task Manager</h1>
          <div className="flex items-center space-x-4">
            <div className="relative group">
              <div className="w-8 h-8 rounded-full bg-blue-400 flex items-center justify-center text-white font-medium cursor-pointer">
                {user?.username?.charAt(0).toUpperCase()}
              </div>
              <div className="absolute -bottom-8 left-1/2 transform -translate-x-1/2 bg-gray-800/80 backdrop-blur-sm text-white px-2 py-1 rounded text-sm whitespace-nowrap opacity-0 group-hover:opacity-100 transition-all duration-200 ease-in-out">
                {user?.username}
              </div>
            </div>
            <button
              onClick={logout}
              className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700"
            >
              Logout
            </button>
          </div>
        </div>
        <div className="mb-6">
          <button
            onClick={() => setShowModal(true)}
            className="flex items-center px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Add New Task
          </button>
        </div>
        <TaskList tasks={tasks} onStatusChange={handleStatusChange} />
      </div>
      {showModal && (
        <CreateTaskModal
          users={users}
          onClose={() => setShowModal(false)}
          onTaskCreated={fetchTasks}
        />
      )}
    </div>
  );
}