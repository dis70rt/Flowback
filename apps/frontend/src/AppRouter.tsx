
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { WorkspaceLayout } from './components/layout/WorkspaceLayout';
import { Workspace } from './screens/Workspace';

export const AppRouter = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<WorkspaceLayout />}>
          <Route index element={<Navigate to="/cases" replace />} />
          <Route path="cases" element={<Workspace />} />
          <Route path="customers" element={<div className="p-8 text-white">Customers View (Coming Soon)</div>} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
};
