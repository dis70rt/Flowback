import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { WorkspaceLayout } from './components/layout/WorkspaceLayout';
import { Workspace } from './screens/Workspace';
import { Customers } from './screens/Customers';

export const AppRouter = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<WorkspaceLayout />}>
          <Route index element={<Navigate to="/cases" replace />} />
          <Route path="cases" element={<Workspace />} />
          <Route path="customers" element={<Customers />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
};
