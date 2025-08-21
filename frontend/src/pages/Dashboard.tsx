import * as React from 'react';
import { useState, useEffect } from 'react';
import { format, startOfWeek, addDays } from 'date-fns';
import { 
  Utensils, 
  Dumbbell, 
  Target, 
  TrendingUp, 
  Calendar,
  Plus,
  Flame,
  Zap
} from 'lucide-react';
import { api } from '../services/api';
import toast from 'react-hot-toast';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

interface DashboardData {
  meals: any[];
  workouts: any[];
  summary: {
    calories: number;
    protein: number;
    carbs: number;
    fat: number;
    total_sets: number;
    total_reps: number;
    total_volume: number;
  };
}

interface WeeklyData {
  days: string[];
  calories: number[];
  protein: number[];
  volume: number[];
  weekly_totals: {
    calories: number;
    protein: number;
    volume: number;
  };
  goals: {
    weekly_calories: number;
    weekly_protein: number;
    weekly_volume: number;
  };
  progress_percent: {
    calories: number;
    protein: number;
    volume: number;
  };
}

const Dashboard: React.FC = () => {
  const [dashboardData, setDashboardData] = useState<DashboardData | null>(null);
  const [weeklyData, setWeeklyData] = useState<WeeklyData | null>(null);
  const [selectedDate, setSelectedDate] = useState(new Date());
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    fetchDashboardData();
    fetchWeeklyData();
  }, [selectedDate]);

  const fetchDashboardData = async () => {
    try {
      const dateStr = format(selectedDate, 'yyyy-MM-dd');
      const response = await api.get(`/dashboard?date=${dateStr}`);
      setDashboardData(response.data);
    } catch (error) {
      toast.error('Failed to fetch dashboard data');
    } finally {
      setIsLoading(false);
    }
  };

  const fetchWeeklyData = async () => {
    try {
      const weekStart = startOfWeek(selectedDate, { weekStartsOn: 1 });
      const startStr = format(weekStart, 'yyyy-MM-dd');
      const response = await api.get(`/dashboard/weekly?start=${startStr}`);
      setWeeklyData(response.data);
    } catch (error) {
      toast.error('Failed to fetch weekly data');
    }
  };

  const getProgressColor = (percentage: number) => {
    if (percentage >= 80) return 'text-success-600';
    if (percentage >= 60) return 'text-warning-600';
    return 'text-danger-600';
  };

  const getProgressBarColor = (percentage: number) => {
    if (percentage >= 80) return 'bg-success-500';
    if (percentage >= 60) return 'bg-warning-500';
    return 'bg-danger-500';
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
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
          <p className="text-gray-600">Track your fitness progress and nutrition</p>
        </div>
        <div className="flex items-center space-x-4">
          <input
            type="date"
            value={format(selectedDate, 'yyyy-MM-dd')}
            onChange={(e) => setSelectedDate(new Date(e.target.value))}
            className="input-field w-auto"
          />
        </div>
      </div>

      {/* Quick Stats */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <div className="card">
          <div className="flex items-center">
            <div className="p-2 bg-primary-100 rounded-lg">
              <Flame className="w-6 h-6 text-primary-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Calories</p>
              <p className="text-2xl font-bold text-gray-900">
                {dashboardData?.summary.calories.toFixed(0) || 0}
              </p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center">
            <div className="p-2 bg-success-100 rounded-lg">
              <Zap className="w-6 h-6 text-success-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Protein (g)</p>
              <p className="text-2xl font-bold text-gray-900">
                {dashboardData?.summary.protein.toFixed(1) || 0}
              </p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center">
            <div className="p-2 bg-warning-100 rounded-lg">
              <Dumbbell className="w-6 h-6 text-warning-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Volume</p>
              <p className="text-2xl font-bold text-gray-900">
                {dashboardData?.summary.total_volume.toFixed(0) || 0}
              </p>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center">
            <div className="p-2 bg-danger-100 rounded-lg">
              <Target className="w-6 h-6 text-danger-600" />
            </div>
            <div className="ml-4">
              <p className="text-sm font-medium text-gray-600">Sets</p>
              <p className="text-2xl font-bold text-gray-900">
                {dashboardData?.summary.total_sets || 0}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Weekly Progress Chart */}
      {weeklyData && (
        <div className="card">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">Weekly Progress</h3>
          <div className="h-80">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={weeklyData.days.map((day, index) => ({
                day: format(new Date(day), 'EEE'),
                calories: weeklyData.calories[index],
                protein: weeklyData.protein[index],
                volume: weeklyData.volume[index],
              }))}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="day" />
                <YAxis />
                <Tooltip />
                <Line type="monotone" dataKey="calories" stroke="#0ea5e9" strokeWidth={2} />
                <Line type="monotone" dataKey="protein" stroke="#22c55e" strokeWidth={2} />
                <Line type="monotone" dataKey="volume" stroke="#f59e0b" strokeWidth={2} />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>
      )}

      {/* Progress vs Goals */}
      {weeklyData && (
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="card">
            <h4 className="font-semibold text-gray-900 mb-4">Calories Progress</h4>
            <div className="space-y-3">
              <div className="flex justify-between text-sm">
                <span>Current: {weeklyData.weekly_totals.calories.toFixed(0)}</span>
                <span>Goal: {weeklyData.goals.weekly_calories.toFixed(0)}</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div
                  className={`h-2 rounded-full ${getProgressBarColor(weeklyData.progress_percent.calories)}`}
                  style={{ width: `${Math.min(weeklyData.progress_percent.calories, 100)}%` }}
                ></div>
              </div>
              <p className={`text-sm font-medium ${getProgressColor(weeklyData.progress_percent.calories)}`}>
                {weeklyData.progress_percent.calories.toFixed(1)}% Complete
              </p>
            </div>
          </div>

          <div className="card">
            <h4 className="font-semibold text-gray-900 mb-4">Protein Progress</h4>
            <div className="space-y-3">
              <div className="flex justify-between text-sm">
                <span>Current: {weeklyData.weekly_totals.protein.toFixed(1)}g</span>
                <span>Goal: {weeklyData.goals.weekly_protein.toFixed(1)}g</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div
                  className={`h-2 rounded-full ${getProgressBarColor(weeklyData.progress_percent.protein)}`}
                  style={{ width: `${Math.min(weeklyData.progress_percent.protein, 100)}%` }}
                ></div>
              </div>
              <p className={`text-sm font-medium ${getProgressColor(weeklyData.progress_percent.protein)}`}>
                {weeklyData.progress_percent.protein.toFixed(1)}% Complete
              </p>
            </div>
          </div>

          <div className="card">
            <h4 className="font-semibold text-gray-900 mb-4">Volume Progress</h4>
            <div className="space-y-3">
              <div className="flex justify-between text-sm">
                <span>Current: {weeklyData.weekly_totals.volume.toFixed(0)}</span>
                <span>Goal: {weeklyData.goals.weekly_volume.toFixed(0)}</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div
                  className={`h-2 rounded-full ${getProgressBarColor(weeklyData.progress_percent.volume)}`}
                  style={{ width: `${Math.min(weeklyData.progress_percent.volume, 100)}%` }}
                ></div>
              </div>
              <p className={`text-sm font-medium ${getProgressColor(weeklyData.progress_percent.volume)}`}>
                {weeklyData.progress_percent.volume.toFixed(1)}% Complete
              </p>
            </div>
          </div>
        </div>
      )}

      {/* Recent Activity */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Recent Meals */}
        <div className="card">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold text-gray-900">Recent Meals</h3>
            <button className="btn-primary flex items-center space-x-2">
              <Plus className="w-4 h-4" />
              <span>Log Meal</span>
            </button>
          </div>
          <div className="space-y-3">
            {dashboardData?.meals.slice(0, 5).map((meal, index) => (
              <div key={index} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div>
                  <p className="font-medium text-gray-900">{meal.description}</p>
                  <p className="text-sm text-gray-600">{meal.calories.toFixed(0)} cal</p>
                </div>
                <div className="text-right">
                  <p className="text-sm text-gray-600">{meal.protein.toFixed(1)}g protein</p>
                  <p className="text-xs text-gray-500">{format(new Date(meal.logged_at), 'HH:mm')}</p>
                </div>
              </div>
            ))}
            {(!dashboardData?.meals || dashboardData.meals.length === 0) && (
              <p className="text-gray-500 text-center py-4">No meals logged today</p>
            )}
          </div>
        </div>

        {/* Recent Workouts */}
        <div className="card">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold text-gray-900">Recent Workouts</h3>
            <button className="btn-primary flex items-center space-x-2">
              <Plus className="w-4 h-4" />
              <span>Log Workout</span>
            </button>
          </div>
          <div className="space-y-3">
            {dashboardData?.workouts.slice(0, 5).map((workout, index) => (
              <div key={index} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div>
                  <p className="font-medium text-gray-900">{workout.exercise}</p>
                  <p className="text-sm text-gray-600">{workout.sets} sets × {workout.reps} reps</p>
                </div>
                <div className="text-right">
                  <p className="text-sm text-gray-600">{workout.weight}kg</p>
                  <p className="text-xs text-gray-500">{format(new Date(workout.logged_at), 'HH:mm')}</p>
                </div>
              </div>
            ))}
            {(!dashboardData?.workouts || dashboardData.workouts.length === 0) && (
              <p className="text-gray-500 text-center py-4">No workouts logged today</p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard; 