import React from 'react';
import { render, fireEvent, screen } from '@testing-library/react';
import { ApiKeyProvider } from '../auth/ApiKey';
import { ApiKeyInput } from './ApiKeyInput';

describe('ApiKeyInput', () => {
  const setup = () => {
    render(
      <ApiKeyProvider>
        <ApiKeyInput />
      </ApiKeyProvider>
    );
    const input = screen.getByLabelText('API Key *');
    const button = screen.getByRole('button', { name: 'Submit' });
    return {
      input,
      button
    };
  };

  it('allows entering an API key', () => {
    const { input } = setup();
    fireEvent.change(input, { target: { value: '12345' } });
    expect(input.value).toBe('12345');
  });

  it('calls saveApiKey on form submit', () => {
    const { input, button } = setup();
    fireEvent.change(input, { target: { value: '12345' } });
    fireEvent.click(button);
    expect(localStorage.getItem('apiKey')).toBe('12345');
  });
});
