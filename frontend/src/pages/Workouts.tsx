import * as React from 'react';
import { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { format } from 'date-fns';
import { Dumbbell, Plus, TrendingUp } from 'lucide-react';
import { api } from '../services/api';
import toast from 'react-hot-toast';

interface WorkoutForm {
  exercise: string;
  sets: number;
  reps: number;
  weight: number;
}

interface Workout {
  exercise: string;
  sets: number;
  reps: number;
  weight: number;
  logged_at: string;
}

const Workouts: React.FC = () => {
  const [workouts, setWorkouts] = useState<Workout[]>([]);
  const [isLogging, setIsLogging] = useState(false);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<WorkoutForm>();

  useEffect(() => {
    fetchWorkouts();
  }, []);

  const fetchWorkouts = async () => {
    try {
      const response = await api.get('/dashboard');
      setWorkouts(response.data.workouts || []);
    } catch (error) {
      toast.error('Failed to fetch workouts');
    }
  };

  const onSubmit = async (data: WorkoutForm) => {
    setIsLogging(true);
    try {
      await api.post('/log-strength', data);
      toast.success('Workout logged successfully!');
      reset();
      fetchWorkouts();
    } catch (error: any) {
      toast.error(error.response?.data || 'Failed to log workout');
    } finally {
      setIsLogging(false);
    }
  };

  const calculateVolume = (sets: number, reps: number, weight: number) => {
    return sets * reps * weight;
  };

  const totalVolume = workouts.reduce(
    (acc, workout) => acc + calculateVolume(workout.sets, workout.reps, workout.weight),
    0
  );

  const totalSets = workouts.reduce((acc, workout) => acc + workout.sets, 0);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Workout Tracking</h1>
        <p className="text-gray-600">Log your strength training sessions and track progress</p>
      </div>

      {/* Summary Stats */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="card text-center">
          <div className="p-3 bg-primary-100 rounded-lg w-fit mx-auto mb-3">
            <Dumbbell className="w-6 h-6 text-primary-600" />
          </div>
          <p className="text-sm text-gray-600">Total Volume</p>
          <p className="text-2xl font-bold text-gray-900">{totalVolume.toFixed(0)}</p>
        </div>
        <div className="card text-center">
          <div className="p-3 bg-success-100 rounded-lg w-fit mx-auto mb-3">
            <TrendingUp className="w-6 h-6 text-success-600" />
          </div>
          <p className="text-sm text-gray-600">Total Sets</p>
          <p className="text-2xl font-bold text-gray-900">{totalSets}</p>
        </div>
        <div className="card text-center">
          <div className="p-3 bg-warning-100 rounded-lg w-fit mx-auto mb-3">
            <Dumbbell className="w-6 h-6 text-warning-600" />
          </div>
          <p className="text-sm text-gray-600">Workouts Today</p>
          <p className="text-2xl font-bold text-gray-900">{workouts.length}</p>
        </div>
      </div>

      {/* Log Workout Form */}
      <div className="card">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Log New Workout</h2>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div>
              <label htmlFor="exercise" className="block text-sm font-medium text-gray-700 mb-2">
                Exercise
              </label>
              <input
                {...register('exercise', {
                  required: 'Exercise name is required',
                })}
                type="text"
                id="exercise"
                className="input-field"
                placeholder="e.g., Bench Press"
              />
              {errors.exercise && (
                <p className="mt-1 text-sm text-danger-600">{errors.exercise.message}</p>
              )}
            </div>

            <div>
              <label htmlFor="sets" className="block text-sm font-medium text-gray-700 mb-2">
                Sets
              </label>
              <input
                {...register('sets', {
                  required: 'Number of sets is required',
                  min: { value: 1, message: 'Sets must be at least 1' },
                  max: { value: 20, message: 'Sets cannot exceed 20' },
                })}
                type="number"
                id="sets"
                className="input-field"
                placeholder="3"
              />
              {errors.sets && (
                <p className="mt-1 text-sm text-danger-600">{errors.sets.message}</p>
              )}
            </div>

            <div>
              <label htmlFor="reps" className="block text-sm font-medium text-gray-700 mb-2">
                Reps
              </label>
              <input
                {...register('reps', {
                  required: 'Number of reps is required',
                  min: { value: 1, message: 'Reps must be at least 1' },
                  max: { value: 100, message: 'Reps cannot exceed 100' },
                })}
                type="number"
                id="reps"
                className="input-field"
                placeholder="10"
              />
              {errors.reps && (
                <p className="mt-1 text-sm text-danger-600">{errors.reps.message}</p>
              )}
            </div>

            <div>
              <label htmlFor="weight" className="block text-sm font-medium text-gray-700 mb-2">
                Weight (kg)
              </label>
              <input
                {...register('weight', {
                  required: 'Weight is required',
                  min: { value: 0, message: 'Weight cannot be negative' },
                  max: { value: 1000, message: 'Weight cannot exceed 1000kg' },
                })}
                type="number"
                step="0.5"
                id="weight"
                className="input-field"
                placeholder="60"
              />
              {errors.weight && (
                <p className="mt-1 text-sm text-danger-600">{errors.weight.message}</p>
              )}
            </div>
          </div>

          <button
            type="submit"
            disabled={isLogging}
            className="btn-primary flex items-center space-x-2"
          >
            {isLogging ? (
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
            ) : (
              <>
                <Plus className="w-4 h-4" />
                <span>Log Workout</span>
              </>
            )}
          </button>
        </form>
      </div>

      {/* Workout History */}
      <div className="card">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Today's Workouts</h2>
        {workouts.length > 0 ? (
          <div className="space-y-3">
            {workouts.map((workout, index) => (
              <div key={index} className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                <div className="flex-1">
                  <p className="font-medium text-gray-900">{workout.exercise}</p>
                  <div className="flex space-x-4 mt-1 text-sm text-gray-600">
                    <span>{workout.sets} sets</span>
                    <span>{workout.reps} reps</span>
                    <span>{workout.weight}kg</span>
                    <span className="font-medium">
                      Volume: {calculateVolume(workout.sets, workout.reps, workout.weight).toFixed(0)}
                    </span>
                  </div>
                </div>
                <div className="text-right">
                  <p className="text-xs text-gray-500">
                    {format(new Date(workout.logged_at), 'HH:mm')}
                  </p>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <p className="text-gray-500 text-center py-8">No workouts logged today</p>
        )}
      </div>
    </div>
  );
};

export default Workouts; 