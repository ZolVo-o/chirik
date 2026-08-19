import React from 'react';
import { render, screen } from '@testing-library/react';
import App from './App';

test('renders the authentication screen', () => {
  render(<App />);
  expect(screen.getByText('Чирик')).toBeInTheDocument();
  expect(screen.getByRole('button', { name: 'Войти' })).toBeInTheDocument();
});
