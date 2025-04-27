import { Product } from "@/types/products.type";

const API_URI = process.env.NEXT_PUBLIC_API_URI || 'http://localhost:8080';

export const getProducts = async (): Promise<Product[]> => {
  const response = await fetch(`${API_URI}/api/products`);
  
  if (!response.ok) {
    throw new Error(`Erreur API: ${response.status}`);
  }
  
  // Fix: Don't consume the response twice
  const data = await response.json();
  console.log(data);
  return data;
};