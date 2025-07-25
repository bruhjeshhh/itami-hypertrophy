import axios from "axios";

const API = axios.create({
  baseURL: "http://localhost:8080",  // Backend URL
});

// Automatically attach token if exists
API.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default API;
