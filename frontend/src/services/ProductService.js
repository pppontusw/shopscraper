const API_URL = process.env.REACT_APP_API_URL;

export const fetchProducts = async (apiKey) => {
  let response;
  try {
    response = await fetch(API_URL + "/products", {
      method: 'GET',
      headers: {
        'X-API-KEY': apiKey,
      },
    });

    if (response.status === 401) {
      throw new Error('Invalid API key'); 
    }
    if (!response.ok) {
      throw new Error(`Failed to fetch products: Server responded with status ${response.status}`);
    }

    return await response.json(); 

  } catch (error) {
    console.error('Error fetching products:', error);
    throw new Error(`Failed to fetch products from ${API_URL}/products`);
  }
};
