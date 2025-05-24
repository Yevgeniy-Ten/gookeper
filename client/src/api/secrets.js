import api from './axios';

export const fetchSecrets = async () => {
  const { data } = await api.get('/secrets');
  return data;
};

export const fetchSecret = async (id) => {
  const { data } = await api.get(`/secrets/${id}`);
  return data;
};

export const createSecret = async (secret) => {
  const { data } = await api.post('/secrets', secret);
  return data;
};

export const deleteSecret = async (id) => {
  await api.delete(`/secrets/${id}`);
};
