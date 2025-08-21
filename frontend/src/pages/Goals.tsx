import * as React from 'react';
import { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { Target, TrendingUp, Calendar } from 'lucide-react';
import { api } from '../services/api';
import toast from 'react-hot-toast';

interface GoalsForm {
  daily_calories: number;
  daily_protein: number;
  weekly_workout_volume: number;
}

interface Goals {
  daily_calories: number;
  daily_protein: number;
  weekly_workout_volume: number;
}

const Goals: React.FC = () => {
  const [goals, setGoals] = useState<Goals | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [isSaving, setIsSaving] = useState(false);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<GoalsForm>();

  useEffect(() => {
    fetchGoals();
  }, []);

  const fetchGoals = async () => {
    try {
      const response = await api.get('/goals');
      setGoals(response.data);
      reset(response.data);
    } catch (error) {
      // Goals not set yet, that's okay
      setIsLoading(false);
    }
  };

  const onSubmit = async (data: GoalsForm) => {
    setIsSaving(true);
    try {
      await api.post('/goals/set', data);
      toast.success('Goals updated successfully!');
      setGoals(data);
    } catch (error: any) {
      toast.error(error.response?.data || 'Failed to update goals');
    } finally {
      setIsSaving(false);
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Goals & Targets</h1>
        <p className="text-gray-600">Set and track your fitness and nutrition goals</p>
      </div>

      {/* Current Goals Display */}
      {goals && (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="card text-center">
            <div className="p-3 bg-primary-100 rounded-lg w-fit mx-auto mb-3">
              <Target className="w-6 h-6 text-primary-600" />
            </div>
            <p className="text-sm text-gray-600">Daily Calories</p>
            <p className="text-2xl font-bold text-gray-900">{goals.daily_calories}</p>
            <p className="text-xs text-gray-500">calories per day</p>
          </div>

          <div className="card text-center">
            <div className="p-3 bg-success-100 rounded-lg w-fit mx-auto mb-3">
              <TrendingUp className="w-6 h-6 text-success-600" />
            </div>
            <p className="text-sm text-gray-600">Daily Protein</p>
            <p className="text-2xl font-bold text-gray-900">{goals.daily_protein}g</p>
            <p className="text-xs text-gray-500">grams per day</p>
          </div>

          <div className="card text-center">
            <div className="p-3 bg-warning-100 rounded-lg w-fit mx-auto mb-3">
              <Calendar className="w-6 h-6 text-warning-600" />
            </div>
            <p className="text-sm text-gray-600">Weekly Volume</p>
            <p className="text-2xl font-bold text-gray-900">{goals.weekly_workout_volume}</p>
            <p className="text-xs text-gray-500">total volume per week</p>
          </div>
        </div>
      )}

      {/* Set/Update Goals Form */}
      <div className="card">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">
          {goals ? 'Update Goals' : 'Set Your Goals'}
        </h2>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div>
              <label htmlFor="daily_calories" className="block text-sm font-medium text-gray-700 mb-2">
                Daily Calorie Target
              </label>
              <input
                {...register('daily_calories', {
                  required: 'Daily calorie target is required',
                  min: { value: 1000, message: 'Calories must be at least 1000' },
                  max: { value: 10000, message: 'Calories cannot exceed 10000' },
                })}
                type="number"
                id="daily_calories"
                className="input-field"
                placeholder="2000"
              />
              {errors.daily_calories && (
                <p className="mt-1 text-sm text-danger-600">{errors.daily_calories.message}</p>
              )}
              <p className="mt-1 text-sm text-gray-500">
                Recommended: 1800-2500 calories for most adults
              </p>
            </div>

            <div>
              <label htmlFor="daily_protein" className="block text-sm font-medium text-gray-700 mb-2">
                Daily Protein Target (g)
              </label>
              <input
                {...register('daily_protein', {
                  required: 'Daily protein target is required',
                  min: { value: 50, message: 'Protein must be at least 50g' },
                  max: { value: 500, message: 'Protein cannot exceed 500g' },
                })}
                type="number"
                step="0.1"
                id="daily_protein"
                className="input-field"
                placeholder="150"
              />
              {errors.daily_protein && (
                <p className="mt-1 text-sm text-danger-600">{errors.daily_protein.message}</p>
              )}
              <p className="mt-1 text-sm text-gray-500">
                Recommended: 1.6-2.2g per kg of body weight
              </p>
            </div>

            <div>
              <label htmlFor="weekly_workout_volume" className="block text-sm font-medium text-gray-700 mb-2">
                Weekly Workout Volume Target
              </label>
              <input
                {...register('weekly_workout_volume', {
                  required: 'Weekly volume target is required',
                  min: { value: 100, message: 'Volume must be at least 100' },
                  max: { value: 100000, message: 'Volume cannot exceed 100,000' },
                })}
                type="number"
                id="weekly_workout_volume"
                className="input-field"
                placeholder="5000"
              />
              {errors.weekly_workout_volume && (
                <p className="mt-1 text-sm text-danger-600">{errors.weekly_workout_volume.message}</p>
              )}
              <p className="mt-1 text-sm text-gray-500">
                Sets × Reps × Weight for all exercises
              </p>
            </div>
          </div>

          <button
            type="submit"
            disabled={isSaving}
            className="btn-primary flex items-center space-x-2"
          >
            {isSaving ? (
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
            ) : (
              <>
                <Target className="w-4 h-4" />
                <span>{goals ? 'Update Goals' : 'Set Goals'}</span>
              </>
            )}
          </button>
        </form>
      </div>

      {/* Tips Section */}
      <div className="card">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">Goal Setting Tips</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <h4 className="font-medium text-gray-900 mb-2">Nutrition Goals</h4>
            <ul className="text-sm text-gray-600 space-y-1">
              <li>• Start with your current maintenance calories</li>
              <li>• Add 200-500 calories for muscle gain</li>
              <li>• Subtract 200-500 calories for fat loss</li>
              <li>• Protein should be 25-30% of total calories</li>
            </ul>
          </div>
          <div>
            <h4 className="font-medium text-gray-900 mb-2">Workout Goals</h4>
            <ul className="text-sm text-gray-600 space-y-1">
              <li>• Aim for 10-20 sets per muscle group per week</li>
              <li>• Progressive overload: increase weight gradually</li>
              <li>• Rest 48-72 hours between training same muscle</li>
              <li>• Track volume to ensure consistent progress</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Goals; 