import * as React from 'react';
import { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { format } from 'date-fns';
import { Utensils, Plus, Search, Trash2 } from 'lucide-react';
import { api } from '../services/api';
import toast from 'react-hot-toast';

interface MealForm {
  description: string;
}

interface Meal {
  description: string;
  calories: number;
  protein: number;
  carbs: number;
  fat: number;
  logged_at: string;
}

const Nutrition: React.FC = () => {
  const [meals, setMeals] = useState<Meal[]>([]);
  const [todayMeals, setTodayMeals] = useState<Meal[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isLogging, setIsLogging] = useState(false);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<MealForm>();

  useEffect(() => {
    fetchMeals();
    fetchTodayMeals();
  }, []);

  const fetchMeals = async () => {
    try {
      const response = await api.get('/meals');
      setMeals(response.data);
    } catch (error) {
      toast.error('Failed to fetch meals');
    }
  };

  const fetchTodayMeals = async () => {
    try {
      const response = await api.get('/meals/today');
      setTodayMeals(response.data.meals || []);
    } catch (error) {
      toast.error('Failed to fetch today\'s meals');
    }
  };

  const onSubmit = async (data: MealForm) => {
    setIsLogging(true);
    try {
      await api.post('/log-calories', data);
      toast.success('Meal logged successfully!');
      reset();
      fetchMeals();
      fetchTodayMeals();
    } catch (error: any) {
      toast.error(error.response?.data || 'Failed to log meal');
    } finally {
      setIsLogging(false);
    }
  };

  const calculateTotals = (mealList: Meal[]) => {
    return mealList.reduce(
      (acc, meal) => ({
        calories: acc.calories + meal.calories,
        protein: acc.protein + meal.protein,
        carbs: acc.carbs + meal.carbs,
        fat: acc.fat + meal.fat,
      }),
      { calories: 0, protein: 0, carbs: 0, fat: 0 }
    );
  };

  const todayTotals = calculateTotals(todayMeals);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Nutrition Tracking</h1>
        <p className="text-gray-600">Log your meals and track your daily nutrition</p>
      </div>

      {/* Today's Summary */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="card text-center">
          <div className="p-3 bg-primary-100 rounded-lg w-fit mx-auto mb-3">
            <Utensils className="w-6 h-6 text-primary-600" />
          </div>
          <p className="text-sm text-gray-600">Total Calories</p>
          <p className="text-2xl font-bold text-gray-900">{todayTotals.calories.toFixed(0)}</p>
        </div>
        <div className="card text-center">
          <div className="p-3 bg-success-100 rounded-lg w-fit mx-auto mb-3">
            <Utensils className="w-6 h-6 text-success-600" />
          </div>
          <p className="text-sm text-gray-600">Protein (g)</p>
          <p className="text-2xl font-bold text-gray-900">{todayTotals.protein.toFixed(1)}</p>
        </div>
        <div className="card text-center">
          <div className="p-3 bg-warning-100 rounded-lg w-fit mx-auto mb-3">
            <Utensils className="w-6 h-6 text-warning-600" />
          </div>
          <p className="text-sm text-gray-600">Carbs (g)</p>
          <p className="text-2xl font-bold text-gray-900">{todayTotals.carbs.toFixed(1)}</p>
        </div>
        <div className="card text-center">
          <div className="p-3 bg-danger-100 rounded-lg w-fit mx-auto mb-3">
            <Utensils className="w-6 h-6 text-danger-600" />
          </div>
          <p className="text-sm text-gray-600">Fat (g)</p>
          <p className="text-2xl font-bold text-gray-900">{todayTotals.fat.toFixed(1)}</p>
        </div>
      </div>

      {/* Log Meal Form */}
      <div className="card">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Log New Meal</h2>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-2">
              What did you eat?
            </label>
            <div className="flex space-x-3">
              <input
                {...register('description', {
                  required: 'Please describe what you ate',
                  minLength: {
                    value: 3,
                    message: 'Description must be at least 3 characters',
                  },
                })}
                type="text"
                id="description"
                className="input-field flex-1"
                placeholder="e.g., grilled chicken breast with brown rice and broccoli"
              />
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
                    <span>Log Meal</span>
                  </>
                )}
              </button>
            </div>
            {errors.description && (
              <p className="mt-1 text-sm text-danger-600">{errors.description.message}</p>
            )}
            <p className="mt-1 text-sm text-gray-500">
              Describe your meal naturally - our AI will calculate the nutrition facts
            </p>
          </div>
        </form>
      </div>

      {/* Today's Meals */}
      <div className="card">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Today's Meals</h2>
        {todayMeals.length > 0 ? (
          <div className="space-y-3">
            {todayMeals.map((meal, index) => (
              <div key={index} className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                <div className="flex-1">
                  <p className="font-medium text-gray-900">{meal.description}</p>
                  <div className="flex space-x-4 mt-1 text-sm text-gray-600">
                    <span>{meal.calories.toFixed(0)} cal</span>
                    <span>{meal.protein.toFixed(1)}g protein</span>
                    <span>{meal.carbs.toFixed(1)}g carbs</span>
                    <span>{meal.fat.toFixed(1)}g fat</span>
                  </div>
                </div>
                <div className="text-right">
                  <p className="text-xs text-gray-500">
                    {format(new Date(meal.logged_at), 'HH:mm')}
                  </p>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <p className="text-gray-500 text-center py-8">No meals logged today</p>
        )}
      </div>

      {/* Meal History */}
      <div className="card">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">Meal History</h2>
        {meals.length > 0 ? (
          <div className="space-y-3">
            {meals.slice(0, 10).map((meal, index) => (
              <div key={index} className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                <div className="flex-1">
                  <p className="font-medium text-gray-900">{meal.description}</p>
                  <div className="flex space-x-4 mt-1 text-sm text-gray-600">
                    <span>{meal.calories.toFixed(0)} cal</span>
                    <span>{meal.protein.toFixed(1)}g protein</span>
                    <span>{meal.carbs.toFixed(1)}g carbs</span>
                    <span>{meal.fat.toFixed(1)}g fat</span>
                  </div>
                </div>
                <div className="text-right">
                  <p className="text-xs text-gray-500">
                    {format(new Date(meal.logged_at), 'MMM dd, HH:mm')}
                  </p>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <p className="text-gray-500 text-center py-8">No meals logged yet</p>
        )}
      </div>
    </div>
  );
};

export default Nutrition; 