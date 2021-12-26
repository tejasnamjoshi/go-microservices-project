import Axios, { AxiosInstance } from "axios";
import { createContext, useContext, useMemo } from "react";

export const AxiosContext = createContext<AxiosInstance>(undefined);

export default function AxiosProvider({
    children,
  }: React.PropsWithChildren<unknown>) {
    const axios = useMemo(() => {
      const axios = Axios.create({
        headers: {
          "Content-Type": "application/json",
        },
        baseURL: 'http://localhost:3002'
      });
  
      axios.interceptors.request.use((config) => {
        const token = localStorage.getItem("token");
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
  
        return config;
      });
  
      return axios;
    }, []);
  
    return (
      <AxiosContext.Provider value={axios}>{children}</AxiosContext.Provider>
    );
  }

  export function useAxios() {
    return useContext(AxiosContext);
  }