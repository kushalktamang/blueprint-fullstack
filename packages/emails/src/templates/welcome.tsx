import {
  Body,
  Button,
  Container,
  Head,
  Heading,
  Hr,
  Html,
  Link,
  Preview,
  Section,
  Tailwind,
  Text,
} from "@react-email/components";

interface WelcomeEmailProps {
  userFirstName: string;
}

function WelcomeEmail({ userFirstName }: WelcomeEmailProps) {
  return (
    <Html>
      <Tailwind>
        <Head />
        <Preview>Welcome to Blueprint Fullstack</Preview>
        <Body className="bg-gray-100 font-sans">
          <Container className="bg-white p-8 rounded-lg shadow-sm my-10 mx-auto max-w-150">
            <Heading className="text-2xl font-bold text-gray-800 mt-4">
              Welcome to Bluestack!
            </Heading>
            <Section>
              <Text className="text-gray-700 text-base">Hi {userFirstName},</Text>
              <Text className="text-gray-700 text-base">Thank you for joining!</Text>
            </Section>
            <Section className="my-8 text-center">
              <Button
                className="bg-orange-600 hover:bg-orange-700 text-white font-medium rounded-md px-6 py-3"
                href="/dashboard"
              >
                Get Started
              </Button>
            </Section>
            <Hr className="border-gray-200 my-6" />
            <Section>
              <Text className="text-gray-600 text-sm">
                If you have any questions, feel free to{" "}
                <Link href="/support" className="text-orange-600 underline">
                  contact our support team
                </Link>
                .
              </Text>
            </Section>
            <Section className="mt-8 text-center">
              <Text className="text-gray-500 text-xs">
                © {new Date().getFullYear()} Moks. All rights reserved.
              </Text>
              <Text className="text-gray-500 text-xs">9 Nayabazar Street, Pokhara, Nepal</Text>
            </Section>
          </Container>
        </Body>
      </Tailwind>
    </Html>
  );
}

WelcomeEmail.PreviewProps = {
  userFirstName: "Moks",
} satisfies WelcomeEmailProps;

export default WelcomeEmail;
