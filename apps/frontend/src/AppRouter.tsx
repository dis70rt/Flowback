import { useState } from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { WorkspaceLayout } from './components/layout/WorkspaceLayout';
import { Workspace } from './screens/Workspace';
import { Customers } from './screens/Customers';
import { Overview } from './screens/Overview';
import { SplashScreen } from './components/SplashScreen';

export const AppRouter = () => {
  const [splashDone, setSplashDone] = useState(false);

  return (
    <>
      {/* Dashboard always rendered at full opacity beneath the splash overlay */}
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<WorkspaceLayout />}>
            <Route index element={<Overview />} />
            <Route path="cases" element={<Workspace />} />
            <Route path="customers" element={<Customers />} />
          </Route>
        </Routes>
      </BrowserRouter>

      {/* Splash overlays on top and fades out — dashboard shows through immediately */}
      {!splashDone && <SplashScreen onComplete={() => setSplashDone(true)} />}
    </>
  );
};
