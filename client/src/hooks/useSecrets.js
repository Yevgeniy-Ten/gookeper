// client/src/hooks/useSecrets.js
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { fetchSecrets, fetchSecret, createSecret, deleteSecret } from '../api/secrets';

// Получить все секреты
export function useSecrets() {
  return useQuery({
    queryKey: ['secrets'],
    queryFn: fetchSecrets,
  });
}

// Получить один секрет
export function useSecret(id) {
  return useQuery({
    queryKey: ['secret', id],
    queryFn: () => fetchSecret(id),
    enabled: !!id,
  });
}

// Создать секрет
export function useCreateSecret() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createSecret,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['secrets'] });
    },
  });
}

// Удалить секрет
export function useDeleteSecret() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteSecret,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['secrets'] });
    },
  });
}