'use client';

import { useEffect, useState } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { Product } from '@/types/products.type';
import { createNewProductsSSE } from '@/services/sse.service';

export function useNewProducts() {
  const [isConnected, setIsConnected] = useState(false);
  const [lastNewProduct, setLastNewProduct] = useState<Product | null>(null);
  const [newProducts, setNewProducts] = useState<Product[]>([]);
  const queryClient = useQueryClient();

  useEffect(() => {
    const sseService = createNewProductsSSE();
    
    sseService.onOpen(() => {
      setIsConnected(true);
    });
    
    sseService.onMessage((product) => {
      setLastNewProduct(product);
      
      setNewProducts(prev => [product, ...prev].slice(0, 10));
      
      queryClient.setQueryData(['products'], (oldData: Product[] | undefined) => {
        if (!oldData) return [product];
        
        const exists = oldData.some(p => 
          p.id === product.id || 
          (product.reference && p.reference === product.reference)
        );
        
        if (exists) {
          return oldData.map(p => {
            if (p.id === product.id || (product.reference && p.reference === product.reference)) {
              return { ...p, ...product };
            }
            return p;
          });
        } else {
          return [product, ...oldData];
        }
      });
    });
    
    sseService.onError(() => {
      setIsConnected(false);
    });
    
    sseService.connect();
    
    return () => {
      sseService.disconnect();
      setIsConnected(false);
    };
  }, [queryClient]);

  return { isConnected, lastNewProduct, newProducts };
}