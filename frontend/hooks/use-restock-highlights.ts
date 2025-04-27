'use client';

import { useEffect, useState } from 'react';
import { Product } from '@/types/products.type';

export function useRestockHighlights(lastEvent: Product | null, highlightDuration = 15000) {
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
        
        // Remove product from recent restocks after specified duration
        setTimeout(() => {
          console.log('Removing restocked product with ID:', productId);
          setRecentRestocks(current => {
            const newState = { ...current };
            delete newState[productId];
            return newState;
          });
        }, highlightDuration);
        
        return updated;
      });
    }
  }, [lastEvent, highlightDuration]);

  return recentRestocks;
}