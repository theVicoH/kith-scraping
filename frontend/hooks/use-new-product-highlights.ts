'use client';

import { useEffect, useState } from 'react';
import { Product } from '@/types/products.type';

export function useNewProductHighlights(lastEvent: Product | null) {
  const [recentNewProducts, setRecentNewProducts] = useState<Record<string, boolean>>({});

  useEffect(() => {
    console.log('Current recentNewProducts state:', recentNewProducts);
  }, [recentNewProducts]);

  useEffect(() => {
    if (lastEvent) {
      console.log('Received new product event for:', lastEvent.title);
      
      setRecentNewProducts(prev => {
        const productId = String(lastEvent.id);
        console.log('Adding new product with ID:', productId);
        
        const updated = { ...prev, [productId]: true };
        
        return updated;
      });
    }
  }, [lastEvent]);

  return recentNewProducts;
}