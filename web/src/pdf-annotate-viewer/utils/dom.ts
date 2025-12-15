export const getPosition = (element: HTMLElement) => {
  const clientRect = element.getBoundingClientRect();
  return {
    left: clientRect.left + document.body.scrollLeft,
    top: clientRect.top + document.body.scrollTop,
  };
}

export const copyTextToClipboard = (text: string) => {
  try {
    const input = document.createElement('input');
    input.setAttribute('readonly', 'readonly');
    input.style.position = 'fixed';
    input.style.top = '0';
    input.style.left = '0';
    input.style.opacity = '0';
    input.value = text;
    document.body.appendChild(input);
    input.setSelectionRange(0, text.length);
    input.select();
    document.execCommand('copy');
    document.body.removeChild(input);
    return true;
  } catch (error) {
    console.warn(error);
    return false;
  }
};
