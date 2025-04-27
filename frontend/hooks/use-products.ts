'use client';

import { useQuery } from '@tanstack/react-query';
import { getProducts } from "@/services/products.service";
import { Product } from '@/types/products.type';

export function useProducts() {
  return useQuery<Product[], Error>({
    queryKey: ['products'],
    queryFn: getProducts,
  });
}