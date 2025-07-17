import { useEffect, useState } from "react";
import API from "../api";
import MealsList from "./MealsList";
import WeeklyChart from "./WeeklyChart";

export default function Dashboard({ onLogout }) {
  const [meals, setMeals] = useState([]);
  const [weeklyData, setWeeklyData] = useState(null);

  useEffect(() => {
    loadData();
  }, []);

  async function loadData() {
    const todayMeals = await API.get("/meals/today");
    setMeals(todayMeals.data);

    const weekly = await API.get("/dashboard/weekly");
    setWeeklyData(weekly.data);
  }

  return (
    <div>
      <button onClick={onLogout}>Logout</button>
      <h2>Today’s Meals</h2>
      <MealsList meals={meals} />
      <h2>Weekly Progress</h2>
      {weeklyData && <WeeklyChart data={weeklyData} />}
    </div>
  );
}
