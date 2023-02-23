/*async function submitForm(body) {
  console.log('START SUBMIT FORM', body);

  try {
    body = await preSubmitFunctions(body);
  } catch (e) {
    console.log('e', e);
    alert(e);
    return;
  }

  body = mergeFormWithMeData(body);

  if (!body) return; // avoid saving bad state
  const info = ui.meInfo as any;
  if (!info) return console.log('no meInfo');

  const date = new Date();
  const unixTimestamp = Math.floor(date.getTime() / 1000);
  setLoading(true);
  try {
    const newBody = {
      ...body,
      alert: undefined,
      new_ticket_time: unixTimestamp,
      extras: {
        ...body?.extras,
        alert: body.alert
      }
    };

    await main.saveProfile(
      config.name === 'about' || config.name === 'wanted' ? { ...newBody } : body
    );
    closeModal(true);
  } catch (e) {
    console.log('e', e);
  }
  setLoading(false);
  props?.ReCallBounties();
}

*/
export {};
