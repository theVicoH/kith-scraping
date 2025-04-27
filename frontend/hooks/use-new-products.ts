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
    console.log('Initializing new products SSE connection...');
    const sseService = createNewProductsSSE();
    
    console.log('SSE service created for new products, attempting to connect...');
    
    sseService.onOpen(() => {
      console.log('New products SSE connection opened successfully');
      setIsConnected(true);
    });
    
    
    sseService.onMessage((product) => {
      console.log('Received new product message:', product);
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
    
    sseService.onError((error) => {
      console.error('New products SSE connection error:', error);
      setIsConnected(false);
    });
    
    sseService.connect();
    
    return () => {
      console.log('Cleaning up new products SSE connection');
      sseService.disconnect();
      setIsConnected(false);
    };
  }, [queryClient]);

  return { isConnected, lastNewProduct, newProducts };
}