'use client';

import { useProducts } from '@/hooks/use-products';
import { Button } from '@/components/button';
import Image from 'next/image';

export function ProductList() {
  const { data: products, isLoading, error } = useProducts();

  if (isLoading) {
    return <div className="flex justify-center p-8">Chargement des produits...</div>;
  }

  if (error) {
    return (
      <div className="flex flex-col items-center p-8">
        <p className="text-destructive mb-4">Erreur lors du chargement des produits</p>
        <Button variant="outline" onClick={() => window.location.reload()}>
          Réessayer
        </Button>
      </div>
    );
  }

  if (!products || products.length === 0) {
    return <div className="text-center p-8">Aucun produit trouvé</div>;
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6 p-6">
      {products.map((product) => (
        <div 
          key={product.id} 
          className="border rounded-lg overflow-hidden shadow-sm hover:shadow-md transition-shadow"
        >
          <div className="relative h-64 w-full bg-muted flex items-center justify-center">
            {product.image_url ? (
              <Image
                src={product.image_url}
                alt={product.title}
                fill
                className="object-cover"
                unoptimized={!product.image_url.startsWith('/')}
              />
            ) : (
              <span className="text-muted-foreground text-sm">Pas d&apos;image</span>
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
      ))}
    </div>
  );
}