import { LineChart, Line, XAxis, YAxis, Tooltip, Legend } from "recharts";

export default function WeeklyChart({ data }) {
  const chartData = data.days.map((day, i) => ({
    day,
    calories: data.calories[i],
    protein: data.protein[i],
    volume: data.volume[i],
  }));

  return (
    <LineChart width={600} height={300} data={chartData}>
      <XAxis dataKey="day" />
      <YAxis />
      <Tooltip />
      <Legend />
      <Line type="monotone" dataKey="calories" stroke="red" />
      <Line type="monotone" dataKey="protein" stroke="blue" />
      <Line type="monotone" dataKey="volume" stroke="green" />
    </LineChart>
  );
}
