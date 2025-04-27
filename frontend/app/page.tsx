import { ProductList } from '@/components/product-list';

export default function Home() {
  return (
    <div className="container mx-auto py-8">
      <h1 className="text-3xl font-bold mb-8 text-center">Kith Monitor</h1>
      <ProductList />
    </div>
  );
}
