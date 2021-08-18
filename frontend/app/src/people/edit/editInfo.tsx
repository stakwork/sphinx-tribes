import React, { useEffect, useState, useRef } from "react";
import { useStores } from "../../store";
import { useObserver } from "mobx-react-lite";
import Form from "../../form";
import ConfirmMe from "./../confirmMe";
import type { MeInfo, MeData } from '../../store/ui'
import { emptyMeInfo } from '../../store/ui'
import { ftuxEditMeSchema, meSchema } from '../../form/schema'
import api from '../../api'
import styled, { css } from "styled-components";
import { getHostIncludingDockerHosts } from "../../host";
import { Button } from "../../sphinxUI";

export default function EditInfo(props: any) {
    const { ui, main } = useStores();

    const [loading, setLoading] = useState(false);
    const [warnBeforeClose, setWarnBeforeClose] = useState(false);
    const [initialFormState, setInitialFormState]: any = useState(null);
    const scrollDiv: any = useRef(null)
    const formRef: any = useRef(null)

    function closeModal(override) {
        // if form state has changed confirm that changes will be lost
        if (!override) {
            if (formHasUnsavedChanges()) {
                setWarnBeforeClose(true)
                return
            }
        }

        setWarnBeforeClose(false)
        ui.setEditMe(false);
        // ui.setMeInfo(null);

        if (props.done) props.done()
    }

    function fullStateCompare(no1, no2, r) {
        let result = r

        function foundChange_(name, a, b) {
            console.log('foundChange', name, a, b)
            result = true
        }

        let widgetSchemas: any = meSchema.find(f => f.name === 'extras')

        if (no1 && no2) {
            Object.keys(no1).forEach((name) => {
                if (result) return
                let current = no1[name]
                let previous = no2[name]

                // if its a new multi widget, this will trigger
                if (!previous) foundChange_(name, current, previous)

                // if extras, we're comparing objects
                if (name === 'extras') {
                    Object.keys(current).forEach((c) => {
                        //extras
                        // get schema to see if single or list widget
                        let thisSchema = widgetSchemas.extras.find(f => f.name === c)
                        const single = thisSchema.single
                        const a = current[c]
                        const b = previous[c]

                        if (single) {
                            // compare single widget (single object)
                            Object.keys(a).forEach(n => {
                                if (!b || !b[n]) {
                                    foundChange_(name, a, b)
                                } else {
                                    // console.log('compare single down', a[n], b[n])
                                    if (a[n] != b[n]) foundChange_(name, a[n], b[n])
                                }
                            })
                        } else {
                            // compare list widget (array of objects)
                            Array.isArray(a) && a.forEach((k, i) => {
                                const akey = a[i]
                                const bkey = b[i]
                                if (!b || !bkey) {
                                    foundChange_(name, akey, bkey)
                                } else {
                                    Object.keys(akey).forEach(n => {
                                        // console.log('compare multi down', akey[n], bkey[n])
                                        if (akey[n] != bkey[n]) foundChange_(name, akey, bkey)
                                    })
                                }
                            })
                        }
                    })
                }
                // if not extras, we're comparing values (string or number)
                else if (current != previous) {
                    foundChange_(name, previous, current)
                }
            })
        }

        return result
    }

    function formHasUnsavedChanges() {
        let result = false
        let currentState = formRef && formRef.current && formRef.current && formRef.current.values
        // compare up
        try {
            result = fullStateCompare(currentState, initialFormState, result)
            // compare down
            result = fullStateCompare(initialFormState, currentState, result)
        } catch (e) {
            console.log('formHasUnsavedChanges error', e)
        }

        return result
    }

    async function testChallenge(chal: string) {
        try {
            const me: MeInfo = await api.get(`poll/${chal}`);
            if (me && me.pubkey) {
                ui.setMeInfo(me);
                ui.setEditMe(true);
            }
        } catch (e) {
            console.log(e);
        }
    }

    useEffect(() => {
        try {
            var urlObject = new URL(window.location.href);
            var params = urlObject.searchParams;
            const chal = params.get("challenge");
            if (chal) {
                testChallenge(chal);
            }
        } catch (e) { }
    }, []);

    async function submitForm(body) {
        console.log('SUBMIT FORM', body);
        const info = ui.meInfo as any;
        if (!info) return console.log("no meInfo");
        setLoading(true);
        try {
            const URL = info.url.startsWith("http") ? info.url : `https://${info.url}`;
            const r = await fetch(URL + "/profile", {
                method: "POST",
                body: JSON.stringify({
                    // use docker host (tribes.sphinx), because relay will post to it
                    host: getHostIncludingDockerHosts(),
                    ...body,
                    price_to_meet: parseInt(body.price_to_meet),
                }),
                headers: {
                    "x-jwt": info.jwt,
                    "Content-Type": "application/json",
                },
            });
            if (!r.ok) {
                setLoading(false);
                return alert("Failed to create profile");
            }

            await main.getPeople('')
            closeModal(true)
        } catch (e) {
            console.log('e', e)
        }
        setLoading(false);

    }

    useEffect(() => {
        // we do this so that we can track changes to the form
        // if modal is closed before saving
        let initialValues: MeData = emptyMeInfo;
        if (ui.meInfo) {
            initialValues.id = ui.meInfo.id || 0
            initialValues.pubkey = ui.meInfo.pubkey
            initialValues.owner_alias = ui.meInfo.alias || ""
            initialValues.photo_url = ui.meInfo.photo_url || ""
            initialValues.price_to_meet = ui.meInfo.price_to_meet || 0
            initialValues.description = ui.meInfo.description || ""
            initialValues.extras = { ...ui.meInfo.extras } || {}
        }
        setInitialFormState(initialValues)
    }, [ui.meInfo])



    return useObserver(() => {
        let verb = "Create";
        if (ui.meInfo && ui.meInfo.id) verb = "Edit";

        let initialValues: MeData = emptyMeInfo;

        if (ui.meInfo) {
            initialValues.id = ui.meInfo.id || 0
            initialValues.pubkey = ui.meInfo.pubkey
            initialValues.owner_alias = ui.meInfo.alias || ""
            initialValues.photo_url = ui.meInfo.photo_url || ""
            initialValues.price_to_meet = ui.meInfo.price_to_meet || 0
            initialValues.description = ui.meInfo.description || ""
            initialValues.extras = ui.meInfo.extras || {}
        }



        return (
            <div style={{ ...props.style }}>
                {warnBeforeClose &&
                    <WarnSave>
                        <div>
                            Save changes?
                        </div>

                        <BWrap>
                            <Button
                                text='Back'
                                loading={props.loading}
                                onClick={() => setWarnBeforeClose(false)}
                            />
                            <Button
                                text='Discard Changes'
                                loading={props.loading}
                                onClick={() => closeModal(true)}
                            />

                            <Button
                                text='Save Changes'
                                loading={props.loading}
                                onClick={() => {
                                    // submit form
                                    if (formRef && formRef.current) {
                                        try {
                                            formRef.current.handleSubmit()
                                        } catch (e) {
                                            console.log('e', e)
                                        }
                                    }
                                }}
                            />
                        </BWrap>
                    </WarnSave>
                }

                <B ref={scrollDiv} hide={warnBeforeClose}>

                    {props.ftux &&
                        <Welcome>

                        </Welcome>

                    }
                    {!ui.meInfo && <ConfirmMe />}
                    {ui.meInfo && (
                        <Form
                            formRef={formRef}
                            paged={true}
                            loading={loading}
                            onSubmit={submitForm}
                            scrollDiv={scrollDiv}
                            schema={ftuxEditMeSchema}
                            initialValues={initialValues}
                            extraHTML={
                                ui.meInfo.verification_signature
                                    ? {
                                        twitter: `<span>Post this to your twitter account to verify:</span><br/><strong>Sphinx Verification: ${ui.meInfo.verification_signature}</strong>`,
                                    }
                                    : {}
                            }
                        />
                    )}
                </B>
            </div>
        );
    });
}


const EnvWithScrollBar = ({ thumbColor, trackBackgroundColor }) => css`
                scrollbar-color: ${thumbColor} ${trackBackgroundColor}; // Firefox support
                scrollbar-width: thin;

                &::-webkit-scrollbar {
                    width: 6px;
                height: 100%;
  }

                &::-webkit-scrollbar-thumb {
                    background - color: ${thumbColor};
                background-clip: content-box;
                border-radius: 5px;
                border: 1px solid ${trackBackgroundColor};
  }

                &::-webkit-scrollbar-corner,
                &::-webkit-scrollbar-track {
                    background - color: ${trackBackgroundColor};
  }
}

                `
interface BProps {
    hide: boolean;
}
const B = styled.div<BProps>`
                    display: ${p => p.hide && 'none'};
                    height:calc(100% - 4px);
                    width: calc(100% - 4px);
                    overflow-y:auto;
                    padding:0 20px;
                    box-sizing:border-box;
                    ${EnvWithScrollBar({
    thumbColor: '#5a606c',
    trackBackgroundColor: 'rgba(0,0,0,0)',
})}
                    `


const BWrap = styled.div`
        display: flex;
        flex-direction:column;
        justify-content: space-evenly;
        align-items:center;
        width:100%;
        height:170px;
        min-height:100px;
        margin-top:20px;
                    `;

const WarnSave = styled.div`
        display: flex;
        flex:1;
        flex-direction:column;
        justify-content: center;
        align-items:center;
        color:#fff
                    `;

const Welcome = styled.div`
        display: flex;
        flex:1;
        flex-direction:column;
        justify-content: center;
        align-items:center;
        
                    `;
