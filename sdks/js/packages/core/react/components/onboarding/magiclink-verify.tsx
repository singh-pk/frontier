'use client';

import { Button, Flex, Link, Text, TextField } from '@raystack/apsara';
import React, {
  ComponentPropsWithRef,
  useCallback,
  useEffect,
  useRef,
  useState
} from 'react';
import { Container } from '~/react/components/Container';
import { Header } from '~/react/components/Header';
import { useFrontier } from '~/react/contexts/FrontierContext';
import { hasWindow } from '~/utils/index';

// @ts-ignore
import styles from './onboarding.module.css';

type MagicLinkVerifyProps = ComponentPropsWithRef<typeof Container> & {
  logo?: React.ReactNode;
  title?: string;
};

export const MagicLinkVerify = ({
  logo,
  title = 'Check your email',
  ...props
}: MagicLinkVerifyProps) => {
  const [loading, setLoading] = useState<boolean>(false);
  const { client, config } = useFrontier();
  const [emailParam, setEmailParam] = useState<string>('');
  const [stateParam, setStateParam] = useState<string>('');
  const [codeParam, setCodeParam] = useState<string>('');
  const [otp, setOTP] = useState<string>('');
  const [submitError, setSubmitError] = useState<string>('');
  const isButtonDisabledRef = useRef(true);

  const handleOTPChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const { value } = event.target;
    isButtonDisabledRef.current = value.length === 0;
    if (submitError.length > 0) setSubmitError('');
    setOTP(value);
  };

  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    const emailParam = params.get('email');
    const stateParam = params.get('state');
    const codeParam = params.get('code');

    emailParam && setEmailParam(emailParam);
    stateParam && setStateParam(stateParam);
    codeParam && setCodeParam(codeParam);
  }, []);

  const OTPVerifyHandler = useCallback(
    async (e: React.FormEvent<HTMLFormElement>) => {
      e.preventDefault();
      setLoading(true);
      try {
        if (!client) return;

        await client.frontierServiceAuthCallback({
          strategy_name: 'mailotp',
          code: otp,
          state: stateParam
        });

        const searchParams = new URLSearchParams(
          hasWindow() ? window.location.search : ``
        );
        const redirectURL =
          searchParams.get('redirect_uri') || searchParams.get('redirectURL');

        // @ts-ignore
        window.location = redirectURL ? redirectURL : window.location.origin;
      } catch (error) {
        console.log(error);
        isButtonDisabledRef.current = true;
        setSubmitError('Please enter a valid verification code');
      } finally {
        setLoading(false);
      }
    },
    [otp]
  );

  return (
    <Container {...props}>
      <Flex direction={'column'} gap="medium">
        <Header logo={logo} title={title} />
        {emailParam && (
          <Text>
            We have sent a temporary login link. Please check your inbox at
            <b> {emailParam}</b>
          </Text>
        )}
      </Flex>

      <form onSubmit={OTPVerifyHandler} className={styles.container80}>
        <Flex
          direction="column"
          gap="extra-small"
          className={styles.optInputContainer}
        >
          <TextField
            data-test-id="enter-code"
            autoFocus
            // @ts-ignore
            size="medium"
            key={'code'}
            placeholder="Enter code"
            onChange={handleOTPChange}
            className={styles.textFieldCode}
          />

          <Text size={2} className={styles.error}>
            {submitError && String(submitError)}
          </Text>
        </Flex>

        <Button
          data-test-id="continue-with-login-code"
          size="medium"
          variant="primary"
          className={styles.container}
          disabled={isButtonDisabledRef.current}
          type="submit"
        >
          {loading ? 'Submitting...' : 'Continue with login code'}
        </Button>
      </form>

      <Link href={config.redirectLogin} data-test-id="back-to-login">
        <Text size={2}>Back to login</Text>
      </Link>
    </Container>
  );
};
