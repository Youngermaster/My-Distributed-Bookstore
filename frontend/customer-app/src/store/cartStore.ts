import { create } from "zustand";
import { persist } from "zustand/middleware";

interface CartState {
  cartId: string | null;
  setCartId: (id: string) => void;
  clearCartId: () => void;
}

export const useCartStore = create<CartState>()(
  persist(
    (set) => ({
      cartId: null,
      setCartId: (id) => set({ cartId: id }),
      clearCartId: () => set({ cartId: null }),
    }),
    {
      name: "cart-storage",
    }
  )
);
