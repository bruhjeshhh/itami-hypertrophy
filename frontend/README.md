# Itami Hypertrophy Frontend

A modern, beautiful React frontend for the Itami Hypertrophy fitness tracking application.

## Features

- 🎨 **Modern UI/UX** - Beautiful, responsive design with Tailwind CSS
- 🔐 **Authentication** - Secure login/register with JWT
- 📊 **Dashboard** - Comprehensive fitness and nutrition overview
- 🍽️ **Nutrition Tracking** - Log meals with AI-powered nutrition calculation
- 💪 **Workout Tracking** - Log strength training sessions
- 🎯 **Goal Setting** - Set and track fitness targets
- 📱 **Responsive Design** - Works perfectly on all devices
- ⚡ **Fast Performance** - Built with Vite and React 18

## Tech Stack

- **Frontend**: React 18 + TypeScript
- **Styling**: Tailwind CSS
- **Build Tool**: Vite
- **Routing**: React Router DOM
- **Forms**: React Hook Form
- **Charts**: Recharts
- **Icons**: Lucide React
- **Notifications**: React Hot Toast

## Quick Start

### Prerequisites

- Node.js 18+ 
- npm or yarn
- Backend server running on port 8080

### Installation

1. **Install dependencies:**
   ```bash
   npm install
   ```

2. **Start development server:**
   ```bash
   npm run dev
   ```

3. **Open your browser:**
   Navigate to `http://localhost:5173`

### Build for Production

```bash
npm run build
```

The built files will be in the `dist/` directory.

## Environment Variables

Create a `.env` file in the frontend directory:

```env
# Backend API URL (defaults to localhost:8080)
VITE_API_URL=http://localhost:8080
```

## API Integration

The frontend communicates with your Go backend through the following endpoints:

- **Authentication**: `/login`, `/register`
- **Nutrition**: `/log-calories`, `/meals`, `/meals/today`
- **Workouts**: `/log-strength`
- **Dashboard**: `/dashboard`, `/dashboard/weekly`
- **Goals**: `/goals`, `/goals/set`

## Project Structure

```
src/
├── components/          # Reusable UI components
├── contexts/           # React contexts (Auth)
├── pages/              # Page components
├── services/           # API services
├── App.tsx            # Main app component
├── main.tsx           # Entry point
└── index.css          # Global styles
```

## Features Overview

### Dashboard
- Daily nutrition and workout summaries
- Weekly progress charts
- Goal progress tracking
- Recent activity feed

### Nutrition Tracking
- Natural language meal logging
- Automatic nutrition calculation via Nutritionix API
- Daily and historical meal tracking
- Macro breakdown (calories, protein, carbs, fat)

### Workout Tracking
- Exercise logging with sets, reps, and weight
- Volume calculation
- Workout history
- Progress tracking

### Goal Management
- Set daily calorie and protein targets
- Set weekly workout volume goals
- Progress visualization
- Helpful tips and recommendations

## Deployment

### Option 1: Static Hosting (Recommended)

1. Build the project:
   ```bash
   npm run build
   ```

2. Deploy the `dist/` folder to any static hosting service:
   - Vercel
   - Netlify
   - GitHub Pages
   - AWS S3 + CloudFront

### Option 2: Docker

1. Create a Dockerfile:
   ```dockerfile
   FROM nginx:alpine
   COPY dist/ /usr/share/nginx/html/
   COPY nginx.conf /etc/nginx/nginx.conf
   EXPOSE 80
   ```

2. Build and run:
   ```bash
   docker build -t itami-frontend .
   docker run -p 80:80 itami-frontend
   ```

## Development

### Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

### Code Style

- TypeScript for type safety
- Tailwind CSS for styling
- React Hook Form for form management
- Consistent component structure

## Troubleshooting

### Common Issues

1. **Backend Connection Error**
   - Ensure your Go backend is running on port 8080
   - Check the proxy configuration in `vite.config.ts`

2. **Build Errors**
   - Clear `node_modules` and reinstall: `rm -rf node_modules && npm install`
   - Check TypeScript errors: `npm run lint`

3. **Styling Issues**
   - Ensure Tailwind CSS is properly configured
   - Check that all CSS classes are imported

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

This project is part of the Itami Hypertrophy application. 