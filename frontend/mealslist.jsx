export default function MealsList({ meals }) {
  if (!meals.length) return <p>No meals logged today</p>;

  return (
    <ul>
      {meals.map((m) => (
        <li key={m.id}>
          {m.description} - {m.calories} kcal | P: {m.protein}g | C: {m.carbs}g | F: {m.fat}g
        </li>
      ))}
    </ul>
  );
}
