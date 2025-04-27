'use client';

import { useProducts } from '@/hooks/use-products';
import { useProductRestocks } from '@/hooks/use-product-restocks';
import { useRestockHighlights } from '@/hooks/use-restock-highlights';
import { useNewProducts } from '@/hooks/use-new-products';
import { useNewProductHighlights } from '@/hooks/use-new-product-highlights';
import { Button } from '@/components/button';
import { ProductCard } from '@/components/product-card';

export function ProductList() {
  const { data: products, isLoading, error } = useProducts();
  const { isConnected: isRestockConnected, lastEvent: lastRestockEvent } = useProductRestocks();
  const { isConnected: isNewProductConnected, lastNewProduct } = useNewProducts();
  const recentRestocks = useRestockHighlights(lastRestockEvent);
  const recentNewProducts = useNewProductHighlights(lastNewProduct);

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
    <>
      <div className="mb-4 flex justify-between items-center">
        <h2 className="text-xl font-semibold">Produits</h2>
        <div className="text-sm flex gap-4">
          {isRestockConnected ? (
            <span className="text-green-600 flex items-center">
              <span className="h-2 w-2 bg-green-600 rounded-full mr-2 animate-pulse"></span>
              Connecté aux notifications de restocks
            </span>
          ) : (
            <span className="text-red-600">Non connecté aux notifications de restocks</span>
          )}
          
          {isNewProductConnected ? (
            <span className="text-blue-600 flex items-center">
              <span className="h-2 w-2 bg-blue-600 rounded-full mr-2 animate-pulse"></span>
              Connecté aux notifications de nouveaux produits
            </span>
          ) : (
            <span className="text-red-600">Non connecté aux notifications de nouveaux produits</span>
          )}
        </div>
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6 p-6">
        {products?.map((product) => {
          const productId = String(product.id);
          const isRecentRestock = recentRestocks[productId] === true;
          const isNewProduct = recentNewProducts[productId] === true;
          
          return (
            <ProductCard 
              key={product.id}
              product={product}
              isRecentRestock={isRecentRestock}
              isNewProduct={isNewProduct}
            />
          );
        })}
      </div>
    </>
  );
}