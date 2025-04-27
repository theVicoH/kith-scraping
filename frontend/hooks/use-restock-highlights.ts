'use client';

import { useEffect, useState } from 'react';
import { Product } from '@/types/products.type';

export function useRestockHighlights(lastEvent: Product | null) {
  const [recentRestocks, setRecentRestocks] = useState<Record<string, boolean>>({});

  useEffect(() => {
    console.log('Current recentRestocks state:', recentRestocks);
  }, [recentRestocks]);

  useEffect(() => {
    if (lastEvent) {
      console.log('Received restock event for:', lastEvent.title);
      
      setRecentRestocks(prev => {
        const productId = String(lastEvent.id);
        console.log('Adding restocked product with ID:', productId);
        
        const updated = { ...prev, [productId]: true };
        
        return updated;
      });
    }
  }, [lastEvent]);

  return recentRestocks;
}