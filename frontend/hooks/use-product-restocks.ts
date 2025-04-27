'use client';

import { useEffect, useState } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { Product } from '@/types/products.type';
import { createProductRestockSSE } from '@/services/sse.service';

export function useProductRestocks() {
  const [isConnected, setIsConnected] = useState(false);
  const [lastEvent, setLastEvent] = useState<Product | null>(null);
  const queryClient = useQueryClient();

  useEffect(() => {
    const sseService = createProductRestockSSE();
    
    sseService.onOpen(() => {
      setIsConnected(true);
    });
    
    sseService.onMessage((product) => {
      setLastEvent(product);
      
      queryClient.setQueryData(['products'], (oldData: Product[] | undefined) => {
        if (!oldData) return oldData;
        
        return oldData.map(p => {
          if (p.id === product.id) {
            return { ...p, in_stock: product.in_stock };
          }
          
          if (product.reference && p.reference === product.reference) {
            return { ...p, in_stock: product.in_stock };
          }
          
          return p;
        });
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

  return { isConnected, lastEvent };
}