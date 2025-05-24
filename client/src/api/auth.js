import api from './axios';

export const login = async (login, password) => {
  const { data } = await api.post('/login', { login, password });
  return data;
};

export const register = async (login, password) => {
  const { data } = await api.post('/register', { login, password });
  return data;
}; 