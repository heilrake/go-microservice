import { getDriverCarsServer, getDriverServer } from "@/entities/driver/server";
import { DriverDataProvider } from "@/entities/providers/DriverDataProvider";

export default async function DriverLayout({ children }: { children: React.ReactNode }) {
   const driver = await getDriverServer();
   
   const cars = driver ? await getDriverCarsServer() : [];

   return (
      <DriverDataProvider
        initialDriver={driver}
        initialCars={cars}
      >
        {children}
      </DriverDataProvider>
    );
 }