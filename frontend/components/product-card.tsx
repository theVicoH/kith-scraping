'use client';

import { Product } from '@/types/products.type';
import { Button } from '@/components/button';
import Image from 'next/image';

interface ProductCardProps {
  product: Product;
  isRecentRestock: boolean;
  isNewProduct: boolean;
}

export function ProductCard({ product, isRecentRestock, isNewProduct }: ProductCardProps) {  
  return (
    <div 
      className={`rounded-lg overflow-hidden transition-all duration-300 ring-2 ${
        isRecentRestock 
          ? 'ring-green-500' 
          : isNewProduct
            ? 'ring-blue-500'
            : 'ring-gray-100'
      }`}
    >
      <div className="relative h-64 w-full bg-muted flex items-center justify-center">
        {product.image_url ? (
          <Image
            src={product.image_url}
            alt={product.title}
            fill
            priority={true}
            className="object-cover"
            unoptimized={!product.image_url.startsWith('/')}
          />
        ) : (
          <span className="text-muted-foreground text-sm">Pas d&apos;image</span>
        )}
        {isRecentRestock && (
          <div className="absolute top-2 right-2 bg-green-500 text-white text-xs px-2 py-1 rounded-full font-bold">
            Restock!
          </div>
        )}
        {isNewProduct && (
          <div className="absolute top-2 right-2 bg-blue-500 text-white text-xs px-2 py-1 rounded-full font-bold">
            Nouveau!
          </div>
        )}
      </div>
      <div className="p-4">
        <h3 className="font-medium truncate">{product.title}</h3>
        <div className="flex justify-between items-center mt-2">
          <span className="font-bold">{(product.price / 100).toFixed(2)} €</span>
          <span className={`px-2 py-1 rounded text-xs ${product.in_stock ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
            {product.in_stock ? 'En stock' : 'Épuisé'}
          </span>
        </div>
        <a 
          href={product.product_url} 
          target="_blank" 
          rel="noopener noreferrer"
          className="mt-3 block"
        >
          <Button variant="outline" className="w-full">
            Voir sur Kith
          </Button>
        </a>
      </div>
    </div>
  );
}